package user

import (
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"pnas/bt"
	"pnas/category"
	"pnas/db"
	"pnas/log"
	"pnas/phttp"
	"pnas/prpc"
	"pnas/setting"
	"pnas/utils"
	"pnas/video"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type UserManger struct {
	UserTorrents
	mtx          sync.Mutex
	users        map[ID]*User
	genHslRecord map[video.ID]bool
	categorySer  category.Service

	cudaQueue utils.TaskQueue
	qsvQueue  utils.TaskQueue
	soQueue   utils.TaskQueue
}

func (um *UserManger) Init() {
	ut := &UserTorrentsImpl{}
	ut.Init()
	um.UserTorrents = ut
	um.mtx.Lock()
	defer um.mtx.Unlock()

	um.users = make(map[ID]*User)
	um.genHslRecord = make(map[video.ID]bool)

	um.categorySer = &category.Manager{}
	um.categorySer.Init()

	um.cudaQueue.Init(utils.WithMaxQueue(1024))
	um.qsvQueue.Init(utils.WithMaxQueue(1024))
	um.soQueue.Init(utils.WithMaxQueue(1024))
}

func (um *UserManger) Close() {
	um.cudaQueue.Close()
	um.qsvQueue.Close()
	um.soQueue.Close()
}

func (um *UserManger) Login(email string, passwd string) (*User, error) {
	sql := "select id from pnas.user where email=? and passwd=?"
	var userId ID
	err := db.QueryRow(sql, email, passwd).Scan(&userId)
	if err != nil {
		return nil, err
	}
	return um.LoadUser(userId)
}

func (um *UserManger) LoadUser(userId ID) (*User, error) {
	um.mtx.Lock()
	user, ok := um.users[userId]
	if ok {
		um.mtx.Unlock()
		return user, nil
	}
	um.mtx.Unlock()

	user, err := loadUser(userId)
	if err != nil {
		return nil, err
	}
	err = um.addUser(user)
	if err != nil {
		return nil, err
	}

	sql := `select t.version, t.info_hash from user_torrent u 
					left join torrent t on u.torrent_id = t.id where u.user_id=?`
	rows, err := db.Query(sql, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var hashs []bt.InfoHash
	for rows.Next() {
		var infoHash bt.InfoHash
		err = rows.Scan(&infoHash.Version, &infoHash.Hash)
		if err != nil {
			log.Warnf("[user] %d query user torrent %s err: %v", userId, hex.EncodeToString([]byte(infoHash.Hash)), err)
			continue
		}
		hashs = append(hashs, infoHash)
	}

	for _, infoHash := range hashs {
		um.AddUserTorrent(userId, infoHash)
	}
	return user, nil
}

func (um *UserManger) addUser(user *User) error {
	um.mtx.Lock()
	defer um.mtx.Unlock()
	um.users[user.userInfo.Id] = user
	return nil
}

func (um *UserManger) getUser(id ID) *User {
	um.mtx.Lock()
	u, ok := um.users[id]
	if !ok {
		um.mtx.Unlock()
		return nil
	}
	um.mtx.Unlock()
	return u
}

func (um *UserManger) ChangeUserName(id ID, name string) error {
	u := um.getUser(id)
	if u == nil {
		return errors.New("not exist")
	}

	return u.ChangeUserName(name)
}

func (um *UserManger) HasVideo(userId ID, vid video.ID) bool {
	if userId == AdminId {
		return true
	}
	// TODO cache
	sql := "select id from category_item where creator=? and type_id=? and resource_path=?"
	var c int
	err := db.QueryRow(sql, userId, prpc.CategoryItem_Video, vid).Scan(&c)
	if err != nil {
		return false
	}
	return c > 0
}

func (um *UserManger) QueryBtVideoMetadata(userId ID, infoHash bt.InfoHash) (map[int]*video.Metadata, error) {
	if !um.HasTorrent(userId, infoHash) {
		return nil, errors.New("not found bt")
	}
	t, err := um.GetTorrent(infoHash)
	if err != nil {
		return nil, err
	}

	files := t.GetFiles()
	ret := make(map[int]*video.Metadata)
	for i, f := range files {
		if (f.FileType & bt.FileVideoType) == 0 {
			continue
		}
		m := f.Meta
		if m == nil || !video.HasVideoStream(m) {
			continue
		}
		m.Format.Filename = f.Name
		ret[i] = m
	}
	return ret, nil
}

func (um *UserManger) NewCategoryItem(userId ID, params *category.NewCategoryParams) error {
	user := um.getUser(userId)
	if user == nil {
		return errors.New("not found user")
	}
	if params.ParentId <= 0 {
		params.ParentId = user.GetHomeDirectoryId()
	}
	_, err := um.categorySer.NewItem(params)
	return err
}

func (um *UserManger) DelCategoryItem(userId ID, itemId category.ID) error {
	return um.categorySer.DelItem(int64(userId), itemId)
}

func (um *UserManger) QueryItem(userId ID, itemId category.ID) (*category.CategoryItem, error) {
	user, err := um.LoadUser(userId)
	if err != nil {
		return nil, err
	}
	item, err := um.categorySer.GetItem(int64(userId), itemId)
	if err != nil {
		return nil, err
	}
	if !item.HasReadAuth(int64(user.userInfo.Id)) && userId != AdminId {
		return nil, errors.New("no auth")
	}
	return item, nil
}

func (um *UserManger) QueryItems(userId ID, parentId category.ID) []*category.CategoryItem {
	user, err := um.LoadUser(userId)
	if err != nil {
		log.Warnf("[user] %d query items err: %v", userId, err)
		return []*category.CategoryItem{}
	}

	item, err := um.categorySer.GetItem(int64(userId), parentId)
	if err != nil {
		log.Warnf("[user] %d query items err: %v", userId, err)
		return []*category.CategoryItem{}
	}
	if !item.HasReadAuth(int64(user.userInfo.Id)) && userId != AdminId {
		return []*category.CategoryItem{}
	}
	items, err := um.categorySer.GetItems(int64(userId), item.GetSubItemIds()...)
	if err != nil {
		log.Warnf("[user] %d load items err: %v", userId, err)
	}
	return items
}

type AddBtVideosParams struct {
	UserId         ID
	CategoryItemId category.ID
	InfoHash       bt.InfoHash
	FileIndexes    []int32
}

func (um *UserManger) findSubtitle(infoHash bt.InfoHash, videoFileIndex int32) int32 {
	t, err := um.GetTorrent(infoHash)
	if err != nil {
		return -1
	}

	files := t.GetFiles()
	videoFileName := utils.GetFileName(files[videoFileIndex].Name)
	for i, f := range files {
		if int32(i) == videoFileIndex {
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

func (um *UserManger) findAudioTrack(infoHash bt.InfoHash, videoFileIndex int32) []string {
	t, err := um.GetTorrent(infoHash)
	if err != nil {
		return []string{}
	}

	baseInfo := t.GetBaseInfo()
	files := t.GetFiles()
	var r []string
	videoFileName := utils.GetFileName(files[videoFileIndex].Name)
	for i, f := range files {
		if int32(i) == videoFileIndex {
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

func (um *UserManger) AddBtVideos(params *AddBtVideosParams) error {

	user, err := um.LoadUser(params.UserId)
	if err != nil {
		return errors.WithStack(err)
	}
	parentItem, err := um.categorySer.GetItem(int64(params.UserId), category.ID(params.CategoryItemId))
	if err != nil {
		return errors.WithStack(err)
	}
	if !parentItem.IsDirectory() ||
		!parentItem.HasWriteAuth(int64(user.userInfo.Id)) ||
		!um.HasTorrent(params.UserId, params.InfoHash) {
		return errors.New("permission denied")
	}

	t, err := um.GetTorrent(params.InfoHash)
	if err != nil {
		return err
	}

	baseInfo := t.GetBaseInfo()
	files := t.GetFiles()
	for _, i := range params.FileIndexes {
		if int(i) >= len(files) {
			return errors.New("invaild params")
		}
	}
	for _, i := range params.FileIndexes {

		absVideoFN := setting.GS.Bt.SavePath + "/" + files[i].Name

		if (files[i].FileType & bt.FileVideoType) == 0 {
			continue
		}

		v, err := video.GetVideoByFileName(files[i].Name)
		needTryGenHls := false
		if err != nil {
			vid, err := video.New(files[i].Name)
			if err != nil {
				log.Warnf("[user] %d add videos err: %v", params.UserId, err)
				continue
			}
			v.Id = vid
			needTryGenHls = true
		} else {
			needTryGenHls = !v.HlsCreated
		}

		newCParams := &category.NewCategoryParams{
			ParentId:     params.CategoryItemId,
			Creator:      int64(params.UserId),
			TypeId:       prpc.CategoryItem_Video,
			Name:         utils.GetFileName(files[i].Name),
			ResourcePath: strconv.FormatInt(int64(v.Id), 10),
			PosterPath:   "",
			Introduce:    "",
			Auth:         utils.NewBitSet(category.AuthMax),
		}
		item, err := um.categorySer.NewItem(newCParams)
		if err != nil {
			return err
		}

		needGenHls := false
		if needTryGenHls {
			um.mtx.Lock()
			creating, ok := um.genHslRecord[v.Id]
			if !ok || !creating {
				um.genHslRecord[v.Id] = true
				needGenHls = true
			}
			um.mtx.Unlock()
		}

		if needGenHls {
			makeAsUndeal := func() {
				um.mtx.Lock()
				um.genHslRecord[v.Id] = false
				um.mtx.Unlock()
			}
			outDir := setting.GS.Server.HlsPath + fmt.Sprintf("/vid_%d", v.Id)
			audioTracksFN := um.findAudioTrack(params.InfoHash, i)

			um.cudaQueue.TryPut(func() {
				err := video.GenHls(
					&video.GenHlsOpts{
						VideoFileName:     absVideoFN,
						AudioFileNames:    audioTracksFN,
						WantedResolutions: CudaSplitEncoderParams,
						OutDir:            outDir,
						Global:            CudaGlobalDecode,
						GlobalVideoParams: CudaH264GlobalVideoParams,
						GlobalAudioParams: GlobalAudioParams,
					})
				if err != nil {
					err = video.GenHls(
						&video.GenHlsOpts{
							VideoFileName:     absVideoFN,
							AudioFileNames:    audioTracksFN,
							WantedResolutions: CudaEncoderParams2,
							OutDir:            outDir,
							Global:            CudaGlobalDecode2,
							GlobalVideoParams: CudaH264GlobalVideoParams,
							GlobalAudioParams: GlobalAudioParams,
						})
				}

				if err == nil {
					video.VideoHasHls(v.Id)
				} else {
					um.qsvQueue.TryPut(func() {
						err := video.GenHls(
							&video.GenHlsOpts{
								VideoFileName:     absVideoFN,
								AudioFileNames:    audioTracksFN,
								WantedResolutions: QsvSplitEncoderParams,
								Global:            QsvGlobalDecode,
								OutDir:            outDir,
								GlobalVideoParams: QsvH264GlobalVideoParams,
								GlobalAudioParams: GlobalAudioParams,
							})
						if err == nil {
							video.VideoHasHls(v.Id)
							return
						}
						um.soQueue.TryPut(func() {
							err := video.GenHls(
								&video.GenHlsOpts{
									VideoFileName:     absVideoFN,
									AudioFileNames:    audioTracksFN,
									WantedResolutions: SoSplitEncoderParams,
									Global:            SoGlobalDecode,
									OutDir:            outDir,
									GlobalVideoParams: SoH264GlobalVideoParams,
									GlobalAudioParams: GlobalAudioParams,
								})
							if err != nil {
								makeAsUndeal()
							}
						})
					})
				}
			})

			um.soQueue.TryPut(func() {
				video.GenSubtitle(&video.GenSubtitleOpts{
					InputFileName: absVideoFN,
					OutDir:        outDir,
					SubtitleName:  utils.GetFileName(absVideoFN),
					Format:        "webvtt",
					Suffix:        "vtt",
				})
			})

			subtitleFileIndex := um.findSubtitle(params.InfoHash, i)
			if subtitleFileIndex != -1 {
				video.GenSubtitle(&video.GenSubtitleOpts{
					InputFileName: baseInfo.SavePath + "/" + files[subtitleFileIndex].Name,
					OutDir:        outDir,
					SubtitleName:  utils.GetFileName(absVideoFN),
					Format:        "webvtt",
					Suffix:        "vtt",
				})
			}

			go func() {
				rfileName := fmt.Sprintf("vid_%d.jpg", v.Id)
				posterFileName := setting.GS.Server.PosterPath + "/" + rfileName
				err := video.GenPoster(&video.GenPosterParams{
					InputFileName:  absVideoFN,
					OutputFileName: posterFileName,
				})
				if err == nil {
					item.UpdatePosterPath(fmt.Sprintf("/vid_%d.jpg", v.Id))
				}
			}()
		} else {
			rfileName := fmt.Sprintf("vid_%d.jpg", v.Id)
			posterFileName := setting.GS.Server.PosterPath + "/" + rfileName
			fstate, err := os.Stat(posterFileName)
			if err == nil && !fstate.IsDir() {
				item.UpdatePosterPath(rfileName)
			}
		}
	}
	return nil
}

func (um *UserManger) RefreshSubtitle(vid video.ID) error {
	videoFileName, err := video.GetVideoFileName(video.ID(vid))
	if err != nil {
		return errors.New("not found")
	}
	var fs []string
	fn := utils.GetFileName(videoFileName)
	baseName := path.Base(videoFileName)
	walkPath := path.Dir(setting.GS.Bt.SavePath + "/" + utils.FileNameFormat(videoFileName))
	filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warnf("refresh subtitle err: %v", err)
			return err
		}
		if info.IsDir() {
			if path != walkPath {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(info.Name(), fn) && baseName != info.Name() {
			fs = append(fs, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, f := range fs {
		video.GenSubtitle(&video.GenSubtitleOpts{
			InputFileName: f,
			OutDir:        setting.GS.Server.HlsPath + fmt.Sprintf("/vid_%d", vid),
			SubtitleName:  utils.GetFileName(f),
			Format:        "webvtt",
			Suffix:        "vtt",
		})
	}
	um.soQueue.TryPut(func() {
		video.GenSubtitle(&video.GenSubtitleOpts{
			InputFileName: videoFileName,
			OutDir:        setting.GS.Server.HlsPath + fmt.Sprintf("/vid_%d", vid),
			SubtitleName:  utils.GetFileName(videoFileName),
			Format:        "webvtt",
			Suffix:        "vtt",
		})
	})
	return nil
}

func (um *UserManger) IsItemShared(sharedItemId category.ID, itemId category.ID) bool {
	_, err := um.categorySer.GetItem(category.AdminId, sharedItemId)
	if err != nil {
		log.Warnf("not found shared item id %d", sharedItemId)
		return false
	}
	var nextParentId = itemId
	for {
		item, err := um.categorySer.GetItem(category.AdminId, nextParentId)
		if err != nil {
			log.Warnf("not found shared item id :%d, next parent: %d, share item id: %d", itemId, nextParentId, sharedItemId)
			return false
		}
		ii := item.GetItemInfo()
		if ii.Id == sharedItemId {
			return true
		}
		nextParentId = ii.ParentId
	}
}

func writeSubtitle2Item(item *category.CategoryItem, rpcSubtitle *prpc.SubtitleFile) error {
	if item.GetType() != prpc.CategoryItem_Video {
		return errors.New("error type")
	}
	dir := video.GetHlsPlayListPath(item.GetVideoId())
	_, err := os.Stat(dir)
	if err != nil {
		return err
	}
	ext := path.Ext(rpcSubtitle.GetName())
	if phttp.IsHtml5SupportSubtitle(ext) {
		fullPath := path.Join(dir, rpcSubtitle.Name)
		fd, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		fd.Write(rpcSubtitle.Content)
		fd.Close()
	} else {
		err := video.GenSubtitle(&video.GenSubtitleOpts{
			SubtitleContent: rpcSubtitle.Content,
			OutDir:          dir,
			SubtitleName:    utils.GetFileName(rpcSubtitle.Name),
			Format:          "webvtt",
			Suffix:          "vtt",
		})
		return err
	}
	return nil
}

func (um *UserManger) UploadSubtitle(userId ID, req *prpc.UploadSubtitleReq) error {
	item, err := um.categorySer.GetItem(int64(userId), category.ID(req.ItemId))
	if err != nil {
		return err
	}
	if item.GetType() == prpc.CategoryItem_Video {
		for _, s := range req.Subtitles {
			writeSubtitle2Item(item, s)
		}
	} else {
		subItems, err := um.categorySer.GetItems(int64(userId), item.GetSubItemIds()...)
		if err != nil {
			return err
		}
		var itemNames []string
		for _, item := range subItems {
			itemNames = append(itemNames, item.GetItemInfo().Name)
		}
		itemEpisodeMap := utils.ParseEpisode(itemNames)
		var subtitleName []string
		for _, s := range req.Subtitles {
			subtitleName = append(subtitleName, s.GetName())
		}
		subtitleEpisodeMap := utils.ParseEpisode(subtitleName)
		for ep, i := range itemEpisodeMap {
			j, ok := subtitleEpisodeMap[ep]
			if !ok {
				continue
			}
			item := subItems[i]
			writeSubtitle2Item(item, req.Subtitles[j])
		}
	}
	return nil
}
