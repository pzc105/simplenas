package user

import (
	"fmt"
	"os"
	"path"
	"pnas/bt"
	"pnas/category"
	"pnas/db"
	"pnas/phttp"
	"pnas/prpc"
	"pnas/ptype"
	"pnas/user/task"
	"pnas/utils"
	"pnas/video"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

type UserManger struct {
	IMagnetSharesService
	bt.UserTorrents
	mtx         sync.Mutex
	users       map[ptype.UserID]*User
	categorySer category.IService

	tasks task.ITasks
}

func (um *UserManger) Init() {
	ut := &bt.UserTorrentsImpl{}
	ut.Init()
	um.UserTorrents = ut
	um.mtx.Lock()
	defer um.mtx.Unlock()

	um.users = make(map[ptype.UserID]*User)

	um.categorySer = &category.Manager{}
	um.categorySer.Init()

	var magnetShares MagnetSharesService
	magnetShares.Init(um.categorySer, um.UserTorrents)
	um.IMagnetSharesService = &magnetShares

	ts := &task.TasksIml{}
	ts.Init()
	um.tasks = ts
}

func (um *UserManger) Close() {
	um.UserTorrents.Close()
}

func (um *UserManger) Login(email string, passwd string) (*User, error) {
	sql := "select id from pnas.user where email=? and passwd=?"
	var userId ptype.UserID
	err := db.QueryRow(sql, email, passwd).Scan(&userId)
	if err != nil {
		return nil, err
	}
	return um.LoadUser(userId)
}

func (um *UserManger) LoadUser(userId ptype.UserID) (*User, error) {
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
	return user, nil
}

func (um *UserManger) addUser(user *User) error {
	um.mtx.Lock()
	defer um.mtx.Unlock()
	um.users[user.userInfo.Id] = user
	return nil
}

func (um *UserManger) getUser(id ptype.UserID) *User {
	u, _ := um.LoadUser(id)
	return u
}

func (um *UserManger) ChangeUserName(id ptype.UserID, name string) error {
	u := um.getUser(id)
	if u == nil {
		return errors.New("not exist")
	}

	return u.ChangeUserName(name)
}

type ChangePasswordParams struct {
	UserId      ptype.UserID
	Email       string
	OldPassword string
	NewPassword string
}

func (um *UserManger) ChangePassword(params *ChangePasswordParams) bool {
	if len(params.Email) == 0 {
		sql := "update pnas.user set passwd=? where id=? and passwd=?"
		r, err := db.Exec(sql, params.NewPassword, params.UserId, params.OldPassword)
		if err != nil {
			return false
		}
		ra, _ := r.RowsAffected()
		return ra > 0
	} else {
		sql := "update pnas.user set passwd=? where email=? and passwd=?"
		r, err := db.Exec(sql, params.NewPassword, params.Email, params.OldPassword)
		if err != nil {
			return false
		}
		ra, _ := r.RowsAffected()
		return ra > 0
	}
}

func (um *UserManger) HasVideo(userId ptype.UserID, vid ptype.VideoID) bool {
	if userId == AdminId {
		return true
	}
	// TODO cache
	sql := "select id from category_items where creator=? and type_id=? and resource_path=?"
	var c int
	err := db.QueryRow(sql, userId, prpc.CategoryItem_Video, vid).Scan(&c)
	if err != nil {
		return false
	}
	return c > 0
}

func (um *UserManger) QueryBtVideoMetadata(userId ptype.UserID, infoHash *bt.InfoHash) (map[int]*video.Metadata, error) {
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

func (um *UserManger) NewCategoryItem(userId ptype.UserID, params *category.NewCategoryParams) error {
	user := um.getUser(userId)
	if user == nil {
		return errors.New("not found user")
	}
	if params.ParentId <= 0 {
		params.ParentId = user.GetHomeDirectoryId()
	}
	_, err := um.categorySer.AddItem(params)
	return err
}

func (um *UserManger) CategoryService() category.IService {
	return um.categorySer
}

func (um *UserManger) DelCategoryItem(userId ptype.UserID, itemId ptype.CategoryID) error {
	item, err := um.categorySer.GetItem(userId, itemId)
	if err != nil {
		return err
	}
	if item.GetType() == prpc.CategoryItem_Home {
		return errors.New("can't delete home")
	}
	err = um.categorySer.DelItem(userId, itemId)
	if userId == ptype.AdminId && item.GetType() == prpc.CategoryItem_Video && err == nil {
		vidStr := item.GetItemBaseInfo().ResourcePath
		vid, err := strconv.ParseInt(vidStr, 10, 64)
		if err == nil {
			video.RemoveVideo(ptype.VideoID(vid))
		}
	}
	return err
}

type AddBtVideosParams struct {
	UserId         ptype.UserID
	CategoryItemId ptype.CategoryID
	InfoHash       *bt.InfoHash
	FileIndexes    []int32
}

func (um *UserManger) AddBtVideos(params *AddBtVideosParams) error {

	user, err := um.LoadUser(params.UserId)
	if err != nil {
		return err
	}
	parentItem, err := um.categorySer.GetItem(params.UserId, ptype.CategoryID(params.CategoryItemId))
	if err != nil {
		return err
	}
	if !parentItem.IsDirectory() ||
		!parentItem.HasWriteAuth(user.userInfo.Id) ||
		!um.HasTorrent(params.UserId, params.InfoHash) {
		return errors.New("permission denied")
	}

	t, err := um.GetTorrent(params.InfoHash)
	if err != nil {
		return err
	}
	files := t.GetFiles()
	for _, i := range params.FileIndexes {
		if int(i) >= len(files) {
			return errors.New("invaild params")
		}
	}

	for _, i := range params.FileIndexes {
		if (files[i].FileType & bt.FileVideoType) == 0 {
			continue
		}
		um.tasks.NewVideoTask(&task.NewVideoTaskParams{
			UserId:      params.UserId,
			ParentId:    params.CategoryItemId,
			CategorySer: um.categorySer,
			Bt:          um,
			InfoHash:    params.InfoHash,
			BtFileIndex: int(i),
		})
	}
	return nil
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

func (um *UserManger) UploadSubtitle(userId ptype.UserID, req *prpc.UploadSubtitleReq) error {
	item, err := um.categorySer.GetItem(userId, ptype.CategoryID(req.ItemId))
	if err != nil {
		return err
	}
	if item.GetType() == prpc.CategoryItem_Video {
		for _, s := range req.Subtitles {
			writeSubtitle2Item(item, s)
		}
	} else {
		subItems, err := um.categorySer.GetItems(userId, item.GetSubItemIds()...)
		if err != nil {
			return err
		}
		var itemNames []string
		for _, item := range subItems {
			itemNames = append(itemNames, item.GetItemBaseInfo().Name)
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

func (um *UserManger) GetTasks() task.ITasks {
	return um.tasks
}

type RenameBtVideoNameParams struct {
	Who      ptype.UserID
	ParentId ptype.CategoryID
	RefName  string
	NumWidth int
}

func (um *UserManger) RenameBtVideoName(params *RenameBtVideoNameParams) error {
	cate := um.categorySer
	pitem, err := cate.GetItem(params.Who, params.ParentId)

	if err != nil {
		return err
	}

	if !pitem.IsDirectory() {
		var refname string
		if len(params.RefName) > 0 {
			refname = params.RefName
		} else {
			ppitem, err := cate.GetItem(params.Who, pitem.GetItemBaseInfo().ParentId)
			if err == nil {
				refname = ppitem.GetName()
			}
		}

		if pitem.GetType() != prpc.CategoryItem_Video {
			return errors.New("not a video")
		}
		sname := pitem.GetName()
		ep, err := utils.ParseEpisode2(sname)
		if err != nil {
			return err
		}
		pitem.Rename(fmt.Sprintf("%s %0*d", refname, params.NumWidth, ep))
		cate.RefreshItem(pitem.GetParentId())
		return nil
	}

	var refname string
	if len(params.RefName) > 0 {
		refname = params.RefName
	} else {
		refname = pitem.GetName()
	}

	sudIds := pitem.GetSubItemIds()
	for _, id := range sudIds {
		item, err := cate.GetItem(params.Who, id)
		if err != nil || item.GetType() != prpc.CategoryItem_Video {
			continue
		}
		sname := item.GetName()
		ep, err := utils.ParseEpisode2(sname)
		if err != nil {
			continue
		}
		item.Rename(fmt.Sprintf("%s %0*d", refname, params.NumWidth, ep))
	}
	cate.RefreshItem(pitem.GetId())
	return nil
}
