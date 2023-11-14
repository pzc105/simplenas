package task

import (
	"errors"
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
)

type videoTask struct {
	RawTask
	parentId    ptype.CategoryID
	userId      ptype.UserID
	ghlsTasks   *gHlsTask
	bt          bt.UserTorrents
	infoHash    *bt.InfoHash
	btFileIndex int
	categorySer category.IService
}

type NewVideoTaskParams struct {
	UserId      ptype.UserID
	ParentId    ptype.CategoryID
	CategorySer category.IService
	Bt          bt.UserTorrents
	InfoHash    *bt.InfoHash
	BtFileIndex int
	Callback    TaskCallback

	mgr      ITasks
	ghlsTask *gHlsTask
}

func newVideoTask(params *NewVideoTaskParams) ITask {
	var vt videoTask
	vt.RawTask.init(params.mgr)
	vt.userId = params.UserId
	vt.parentId = params.ParentId
	vt.ghlsTasks = params.ghlsTask
	vt.categorySer = params.CategorySer
	vt.bt = params.Bt
	vt.infoHash = params.InfoHash
	vt.btFileIndex = params.BtFileIndex
	vt.callback = params.Callback
	return &vt
}

func (bd *videoTask) hlsCallback(err error) {
	if err == nil || err == ErrExistedVideo {
		bd.into(TaskStatusFinished, nil)
		return
	}
	bd.into(TaskStatusFailed, err)
}

func (v *videoTask) Start() {
	t, err := v.bt.GetTorrent(v.infoHash)
	if err != nil {
		v.into(TaskStatusFailed, err)
		return
	}
	savePath := t.GetBaseInfo().SavePath
	files := t.GetFiles()
	if v.btFileIndex >= len(files) {
		v.into(TaskStatusFailed, errors.New("out of range"))
		return
	}
	f := files[v.btFileIndex]
	fullName := path.Join(savePath, f.Name)
	if !video.IsVideo2(fullName) {
		v.into(TaskStatusFailed, errors.New("is not a video file"))
		return
	}
	vid, err := v.ghlsTasks.add(&AddHlsTaskParams{
		VideoFullName: fullName,
		Callback:      v.hlsCallback,
		AudioTracksFN: v.findAudioTrack(),
	})
	if err != nil && err != ErrExistedVideo {
		v.into(TaskStatusFailed, err)
		return
	}
	outDir := setting.GS().Server.HlsPath + fmt.Sprintf("/vid_%d", vid)

	newCParams := &category.NewCategoryParams{
		ParentId:     v.parentId,
		Creator:      v.userId,
		TypeId:       prpc.CategoryItem_Video,
		Name:         f.Name,
		ResourcePath: strconv.FormatInt(int64(vid), 10),
		Auth:         utils.NewBitSet(category.AuthMax),
	}
	item, err := v.categorySer.AddItem(newCParams)

	if err != nil {
		v.into(TaskStatusFailed, err)
		return
	}
	go func() {
		subtitleFileIndex := v.findSubtitle()
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

func (v *videoTask) findAudioTrack() []string {
	t, err := v.bt.GetTorrent(v.infoHash)
	if err != nil {
		return []string{}
	}

	baseInfo := t.GetBaseInfo()
	files := t.GetFiles()
	var r []string
	videoFileName := utils.GetFileName(files[v.btFileIndex].Name)
	for i, f := range files {
		if i == v.btFileIndex {
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

func (v *videoTask) findSubtitle() int32 {
	t, err := v.bt.GetTorrent(v.infoHash)
	if err != nil {
		return -1
	}

	files := t.GetFiles()
	videoFileName := utils.GetFileName(files[v.btFileIndex].Name)
	for i, f := range files {
		if i == v.btFileIndex {
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
