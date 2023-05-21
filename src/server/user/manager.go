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
	mtx          sync.Mutex
	users        map[ID]*User
	torrents     map[bt.InfoHash]*bt.Torrent
	dlUser       map[bt.InfoHash]map[ID]bool
	genHslRecord map[video.ID]bool

	cudaQueue utils.TaskQueue
	qsvQueue  utils.TaskQueue
	soQueue   utils.TaskQueue
}

func (um *UserManger) Init() {
	um.mtx.Lock()
	defer um.mtx.Unlock()

	um.users = make(map[ID]*User)
	um.torrents = make(map[bt.InfoHash]*bt.Torrent)
	um.dlUser = make(map[bt.InfoHash]map[ID]bool)
	um.genHslRecord = make(map[video.ID]bool)

	um.cudaQueue.Init(utils.WithMaxQueue(1024))
	um.qsvQueue.Init(utils.WithMaxQueue(1024))
	um.soQueue.Init(utils.WithMaxQueue(1024))
}

func (um *UserManger) Close() {
	um.cudaQueue.Close()
	um.qsvQueue.Close()
	um.soQueue.Close()
}

func (um *UserManger) LoadTorrent(infoHash bt.InfoHash) ([]byte, error) {
	sql := `select resume_data from torrent where info_hash=? and version=?`
	var resumeData []byte
	err := db.QueryRow(sql, infoHash.Hash, infoHash.Version).Scan(&resumeData)
	return resumeData, err
}

func (um *UserManger) LoadDownloadingTorrent() [][]byte {
	sql := `select resume_data from user_torrent u 
					left join torrent t on u.torrent_id = t.id`

	rows, err := db.Query(sql)
	if err != nil {
		log.Warnf("[bt] failed to load downloading torrent err: %v", err)
		return [][]byte{}
	}
	defer rows.Close()

	var resumData [][]byte
	for rows.Next() {
		var resume []byte
		err = rows.Scan(&resume)
		if err != nil {
			log.Warnf("[bt] failed to load downloading torrent err: %v", err)
			continue
		}
		resumData = append(resumData, resume)
	}
	return resumData
}

func (um *UserManger) GetAllTorrents() []*bt.Torrent {
	um.mtx.Lock()
	defer um.mtx.Unlock()
	var ret []*bt.Torrent
	for _, v := range um.torrents {
		ret = append(ret, v)
	}
	return ret
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

	user, err := LoadUser(userId)
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

	um.mtx.Lock()
	defer um.mtx.Unlock()
	for _, infoHash := range hashs {
		_, ok := um.dlUser[infoHash]
		if !ok {
			um.dlUser[infoHash] = make(map[ID]bool)
		}
		um.dlUser[infoHash][userId] = true
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

// TODO move this function to bt package
func (um *UserManger) NeedUpdateTorrent(info_hash bt.InfoHash) bool {
	um.mtx.Lock()
	defer um.mtx.Unlock()
	t, ok := um.torrents[info_hash]
	if !ok {
		return true
	}
	return t.GetState() != prpc.BtStateEnum_seeding
}

func updateBtFileType(st *bt.Torrent, fileIndex int, absFileName string) {
	if video.IsVideo(absFileName) {
		st.UpdateFileType(fileIndex, bt.FileVideoType)
		meta, _ := video.GetMetadata(absFileName)
		st.UpdateVideoFileMeta(fileIndex, meta)
	}
	if video.IsSubTitle(absFileName) {
		st.UpdateFileType(fileIndex, bt.FileSubtitleType)
	}
	if video.IsAudio(absFileName) {
		st.UpdateFileType(fileIndex, bt.FileAudioType)
	}
}

func (um *UserManger) UpdateTorrentState(infoHash bt.InfoHash, state prpc.BtStateEnum) error {
	srcState := prpc.BtStateEnum_unknown
	um.mtx.Lock()
	st, ok := um.torrents[infoHash]
	if !ok {
		um.mtx.Unlock()
		return errors.New("not found")
	}
	um.mtx.Unlock()
	srcState = st.GetState()

	if srcState == state {
		return nil
	}

	st.UpdateState(state)

	sql := "update torrent state=? where info_hash=? and version=?"
	db.Exec(sql, state, infoHash.Hash, infoHash.Version)

	if state != prpc.BtStateEnum_seeding {
		return nil
	}

	_, err := um.LoadUser(AdminId)
	if err != nil {
		log.Warnf("[user] %d load err: %v", AdminId, err)
		return nil
	}
	files := st.GetFiles()
	var fileIndexes []int32
	for i := range files {
		fileName := files[i].Name

		if files[i].FileType == bt.FileUnknownType {
			absFileName := setting.GS.Bt.SavePath + "/" + fileName
			updateBtFileType(st, i, absFileName)
		}

		v, err := video.GetVideoByFileName(fileName)
		if err == nil && v.HlsCreated {
			continue
		}

		if ft, _ := st.GetFileType(i); (ft | bt.FileVideoType) != 0 {
			fileIndexes = append(fileIndexes, int32(i))
		}
	}
	if len(fileIndexes) == 0 {
		return nil
	}
	return nil
}

func (um *UserManger) UpdateTorrent(base *bt.TorrentBase, btst prpc.BtStateEnum, fileNames []bt.File, resumeData []byte) {
	um.mtx.Lock()
	st, ok := um.torrents[base.InfoHash]
	if !ok {
		st = bt.NewTorrent(base)
		um.torrents[base.InfoHash] = st
	}
	um.mtx.Unlock()
	st.UpdateBaseInfo(base)
	st.UpdateFilesInfo(fileNames)
	st.UpdateResumeData(resumeData)

	baseInfo := st.GetBaseInfo()
	sql := "update torrent set name=?, state=?, total_size=?, piece_length=?, num_pieces=?, resume_data=? where info_hash=? and version=?"
	_, err := db.Exec(sql, baseInfo.Name, st.GetState(),
		baseInfo.TotalSize, baseInfo.PieceLength, baseInfo.NumPieces,
		st.GetResumeData(), baseInfo.InfoHash.Hash, baseInfo.InfoHash.Version)
	if err != nil {
		log.Warnf("[bt] failed to update torrent %s err: %v", hex.EncodeToString([]byte(base.InfoHash.Hash)), err)
	}
}

type FileCompleted struct {
	InfoHash  bt.InfoHash
	FileIndex int32
}

func (um *UserManger) BtFileStateComplete(fs *FileCompleted) {
	um.mtx.Lock()
	t, ok := um.torrents[fs.InfoHash]
	um.mtx.Unlock()
	if !ok {
		return
	}
	tfs, err := t.GetFileState(int(fs.FileIndex))
	if err != nil {
		log.Warnf("[bt] failed to update torrent %s fileindex %d state err: %v", hex.EncodeToString([]byte(fs.InfoHash.Hash)), fs.FileIndex, err)
		return
	}
	if tfs == prpc.BtFile_completed {
		return
	}

	t.UpdateFileState(int(fs.FileIndex), prpc.BtFile_completed)

	baseInfo := t.GetBaseInfo()
	files := t.GetFiles()

	log.Infof("[bt] torrent: %s file: %s completed", baseInfo.Name, files[fs.FileIndex].Name)

	fileName := files[fs.FileIndex].Name
	absFileName := baseInfo.SavePath + "/" + fileName
	updateBtFileType(t, int(fs.FileIndex), absFileName)

	addVideo := func() {
		admin, err := um.LoadUser(AdminId)
		if err != nil {
			return
		}

		v, err := video.GetVideoByFileName(fileName)
		if err == nil && v.HlsCreated {
			return
		}
		if ft, _ := t.GetFileType(int(fs.FileIndex)); (ft | bt.FileVideoType) == 0 {
			return
		}
		um.AddBtVideos(&AddBtVideosParams{
			UserId:         admin.userInfo.Id,
			CategoryItemId: 2,
			InfoHash:       baseInfo.InfoHash,
			FileIndexes:    []int32{fs.FileIndex},
		})
	}
	addVideo()
}

func (um *UserManger) AddTorrent(userId ID, t *bt.TorrentBase) error {
	um.mtx.Lock()
	_, ok := um.torrents[t.InfoHash]
	if !ok {
		um.torrents[t.InfoHash] = bt.NewTorrent(t)
	}
	_, ok = um.dlUser[t.InfoHash]
	if !ok {
		um.dlUser[t.InfoHash] = make(map[ID]bool)
	}
	_, ok = um.dlUser[t.InfoHash][userId]
	if ok {
		um.mtx.Unlock()
		return errors.New("repeat download")
	}
	um.mtx.Unlock()

	_, err := db.Query("call new_torrent(?, ?, ?, @torrent_id);", t.InfoHash.Version, t.InfoHash.Hash, userId)
	if err != nil {
		log.Warnf("[user] %d failed to add torrent %s err: %v", userId, hex.EncodeToString([]byte(t.InfoHash.Hash)), err)
		return err
	}

	um.mtx.Lock()
	um.dlUser[t.InfoHash][userId] = true
	um.mtx.Unlock()
	return nil
}

func (um *UserManger) hasTorrent(userId ID, infoHash bt.InfoHash) bool {
	um.mtx.Lock()
	defer um.mtx.Unlock()
	i, ok := um.dlUser[infoHash]
	if !ok {
		return false
	}
	_, ok = i[userId]
	return ok
}

func (um *UserManger) HasTorrent(userId ID, infoHash bt.InfoHash) bool {
	if userId == AdminId {
		return true
	}
	return um.hasTorrent(userId, infoHash)
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

func (um *UserManger) RemoveUserTorrent(userId ID, infoHash bt.InfoHash) error {
	if !um.hasTorrent(userId, infoHash) {
		return errors.New("not found")
	}
	sql := "select id from torrent where info_hash=? and version=?"
	rows, err := db.Query(sql, infoHash.Hash, infoHash.Version)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var tid int64
		err := rows.Scan(&tid)
		if err != nil {
			return err
		}
		sql := "delete from user_torrent where user_id=? and torrent_id=?"
		_, err = db.Exec(sql, userId, tid)
		if err != nil {
			return err
		}
	}
	um.mtx.Lock()
	delete(um.dlUser[infoHash], userId)
	um.mtx.Unlock()
	return nil
}

func (um *UserManger) RemoveTorrent(userId ID, infoHash bt.InfoHash) error {
	if !um.HasTorrent(userId, infoHash) {
		return errors.New("not found")
	}
	var uids []string
	um.mtx.Lock()
	us, ok := um.dlUser[infoHash]
	if ok {
		for u, k := range us {
			if k {
				uids = append(uids, fmt.Sprintf("user_id=%d", u))
			}
		}
	}
	um.mtx.Unlock()
	if len(uids) == 0 {
		return nil
	}

	sql := "select id from torrent where info_hash=? and version=?"
	rows, err := db.Query(sql, infoHash.Hash, infoHash.Version)
	if err != nil {
		return err
	}
	defer rows.Close()
	userIdsCond := strings.Join(uids, " or ")
	sql = fmt.Sprintf("delete from user_torrent where %s and torrent_id=?", userIdsCond)
	for rows.Next() {
		var tid int64
		err := rows.Scan(&tid)
		if err != nil {
			return err
		}
		_, err = db.Exec(sql, tid)
		if err != nil {
			return err
		}
	}
	um.mtx.Lock()
	delete(um.dlUser, infoHash)
	um.mtx.Unlock()
	return nil
}

func (um *UserManger) QueryBtVideoMetadata(userId ID, infoHash bt.InfoHash) (map[int]*video.Metadata, error) {
	if !um.HasTorrent(userId, infoHash) {
		return nil, errors.New("not found bt")
	}
	um.mtx.Lock()
	t, ok := um.torrents[infoHash]
	if !ok {
		um.mtx.Unlock()
		return nil, errors.New("not found bt")
	}
	um.mtx.Unlock()

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
	_, err := category.NewItem(params)
	return err
}

func (um *UserManger) DelCategoryItem(userId ID, itemId category.ID) error {
	item, err := category.LoadItem(itemId)
	if err != nil {
		return err
	}
	if !item.HasWriteAuth(int64(userId)) && userId != AdminId {
		return errors.New("not auth")
	}
	return category.RemoveItem(itemId)
}

func (um *UserManger) QueryItem(userId ID, itemId category.ID) (*category.CategoryItem, error) {
	user, err := um.LoadUser(userId)
	if err != nil {
		return nil, err
	}
	item, err := category.LoadItem(itemId)
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

	item, err := category.LoadItem(parentId)
	if err != nil {
		log.Warnf("[user] %d query items err: %v", userId, err)
		return []*category.CategoryItem{}
	}
	if !item.HasReadAuth(int64(user.userInfo.Id)) && userId != AdminId {
		return []*category.CategoryItem{}
	}
	items, err := category.LoadItems(item.GetSubItemIds()...)
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
	um.mtx.Lock()
	t := um.torrents[infoHash]
	um.mtx.Unlock()

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
	um.mtx.Lock()
	t := um.torrents[infoHash]
	um.mtx.Unlock()

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
	parentItem, err := category.LoadItem(category.ID(params.CategoryItemId))
	if err != nil {
		return errors.WithStack(err)
	}
	if !parentItem.IsDirectory() ||
		!parentItem.HasWriteAuth(int64(user.userInfo.Id)) ||
		!um.HasTorrent(params.UserId, params.InfoHash) {
		return errors.New("permission denied")
	}

	um.mtx.Lock()
	t, ok := um.torrents[params.InfoHash]
	um.mtx.Unlock()

	if !ok {
		return errors.New("internal error")
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
		item, err := category.NewItem(newCParams)
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
