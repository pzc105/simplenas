package task

import (
	"fmt"
	"path"
	"pnas/bt"
	"pnas/category"
	"pnas/prpc"
	"pnas/ptype"
	"pnas/setting"
	"pnas/utils"
	"pnas/video"
	"strconv"
	"sync/atomic"
)

type BtHlsTask struct {
	RawTask
	magnetUri    string
	userId       ptype.UserID
	parentId     ptype.CategoryID
	tid          ptype.TorrentID
	st           prpc.BtStateEnum
	bt           bt.UserTorrents
	infoHash     *bt.InfoHash
	ghlsTask     *gHlsTask
	hlsTaskCount atomic.Int32
	categorySer  category.IService
}

type newBtHlsTaskParams struct {
	userId      ptype.UserID
	parentId    ptype.CategoryID
	tasks       *TasksIml
	callback    TaskCallback
	bt          bt.UserTorrents
	ghlsTask    *gHlsTask
	categorySer category.IService
}

func newBtHlsTask(params *newBtHlsTaskParams) ITask {
	var bd BtHlsTask
	bd.userId = params.userId
	bd.parentId = params.parentId
	bd.id = params.tasks.newTaskId()
	bd.status.Store(TaskStatusIniting)
	bd.bt = params.bt
	bd.ghlsTask = params.ghlsTask
	bd.categorySer = params.categorySer
	return &bd
}

func (bd *BtHlsTask) hlsCallback(error) {
	bd.hlsTaskCount.Add(-1)
	if bd.hlsTaskCount.Load() == 0 {
		bd.into(TaskStatusFinished, nil)
	}
}

func (bd *BtHlsTask) Start() {
	res, err := bd.bt.Download(&bt.DownloadParams{
		UserId: bd.userId,
		Req: &prpc.DownloadRequest{
			Type:     prpc.DownloadRequest_MagnetUri,
			Content:  []byte(bd.magnetUri),
			SavePath: setting.GS().Bt.SavePath,
		},
	})
	if err != nil {
		bd.into(TaskStatusInitFailed, err)
		return
	}
	infoHash := bt.TranInfoHash(res.InfoHash)
	t, err := bd.bt.GetTorrent(infoHash)
	if err != nil {
		bd.into(TaskStatusInitFailed, err)
	}
	bd.infoHash = infoHash
	bd.tid = t.GetBaseInfo().Id
	bd.bt.SetTaskCallback(&bt.SetTaskCallbackParams{
		UserId:    bd.userId,
		TaskId:    bd.id,
		TorrentId: bd.tid,
		Callback:  bd.onBtStatus,
	})
}

func (bd *BtHlsTask) onBtStatus(err error, st *prpc.TorrentStatus) {
	if err != nil {
		bd.into(TaskStatusInitFailed, err)
		bd.bt.SetTaskCallback(&bt.SetTaskCallbackParams{
			UserId:    bd.userId,
			TaskId:    bd.id,
			TorrentId: bd.tid,
			Callback:  nil,
		})
		return
	}
	if !bt.IsDownloadAll(bd.st) && bt.IsDownloadAll(st.State) {
		bd.bt.SetTaskCallback(&bt.SetTaskCallbackParams{
			UserId:    bd.userId,
			TaskId:    bd.id,
			TorrentId: bd.tid,
			Callback:  nil,
		})

		t, err := bd.bt.GetTorrent(bd.infoHash)
		if err != nil {
			bd.into(TaskStatusInitFailed, err)
			return
		}
		savePath := t.GetBaseInfo().SavePath
		files := t.GetFiles()
		for i, f := range files {
			index := i
			fullName := path.Join(savePath, f.Name)
			if !video.IsVideo2(fullName) {
				continue
			}
			vid, err := bd.ghlsTask.add(&addHlsTaskParams{
				videoFullName: fullName,
				callback:      bd.hlsCallback,
				audioTracksFN: bd.findAudioTrack(index),
			})
			if err != nil {
				continue
			}
			outDir := setting.GS().Server.HlsPath + fmt.Sprintf("/vid_%d", vid)
			bd.hlsTaskCount.Add(1)

			newCParams := &category.NewCategoryParams{
				ParentId:     bd.parentId,
				Creator:      bd.userId,
				TypeId:       prpc.CategoryItem_Video,
				Name:         f.Name,
				ResourcePath: strconv.FormatInt(int64(vid), 10),
				Auth:         utils.NewBitSet(category.AuthMax),
			}
			item, err := bd.categorySer.AddItem(newCParams)

			if err != nil {
				bd.into(TaskStatusInitFailed, err)
				return
			}
			go func() {
				subtitleFileIndex := bd.findSubtitle(index)
				if subtitleFileIndex != -1 {
					video.GenSubtitle(&video.GenSubtitleOpts{
						InputFileName: path.Join(savePath, files[subtitleFileIndex].Name),
						OutDir:        outDir,
						SubtitleName:  utils.GetFileName(fullName),
						Format:        "webvtt",
						Suffix:        "vtt",
					})
				}

				rfileName := fmt.Sprintf("vid_%d.jpg", vid)
				posterFileName := setting.GS().Server.PosterPath + "/" + rfileName
				err := video.GenPoster(&video.GenPosterParams{
					InputFileName:  fullName,
					OutputFileName: posterFileName,
				})
				if err == nil {
					item.UpdatePosterPath(fmt.Sprintf("/vid_%d.jpg", vid))
				}
			}()
		}
	}
	bd.st = st.State
}

func (bd *BtHlsTask) findAudioTrack(videoFileIndex int) []string {
	t, err := bd.bt.GetTorrent(bd.infoHash)
	if err != nil {
		return []string{}
	}

	baseInfo := t.GetBaseInfo()
	files := t.GetFiles()
	var r []string
	videoFileName := utils.GetFileName(files[videoFileIndex].Name)
	for i, f := range files {
		if i == videoFileIndex {
			continue
		}
		if videoFileName != utils.GetFileName(f.Name) {
			continue
		}
		if (f.FileType & bt.FileAudioType) != 0 {
			r = append(r, baseInfo.SavePath+"/"+files[i].Name)
		}
	}
	return r
}

func (bd *BtHlsTask) findSubtitle(videoFileIndex int) int32 {
	t, err := bd.bt.GetTorrent(bd.infoHash)
	if err != nil {
		return -1
	}

	files := t.GetFiles()
	videoFileName := utils.GetFileName(files[videoFileIndex].Name)
	for i, f := range files {
		if i == videoFileIndex {
			continue
		}
		if videoFileName != utils.GetFileName(f.Name) {
			continue
		}
		if (f.FileType & bt.FileSubtitleType) != 0 {
			return int32(i)
		}
	}
	return -1
}
