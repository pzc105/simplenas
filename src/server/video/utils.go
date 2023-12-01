package video

import (
	"fmt"
	"os"
	"pnas/db"
	"pnas/ptype"
	"pnas/setting"
)

type Video struct {
	Id         ptype.VideoID
	HlsCreated bool
	FileName   string
}

func New(fileName string) (ptype.VideoID, error) {
	var newId ptype.VideoID
	err := db.QueryRow("call new_video(?, @new_video_id)", fileName).Scan(&newId)
	if err != nil {
		return -1, err
	}
	return newId, nil
}

func VideoHasHls(vid ptype.VideoID) error {
	sqlStr := "update video set hls_created=1 where id=?"
	_, err := db.Exec(sqlStr, vid)
	return err
}

func GetVideoByFileName(fileName string) (Video, error) {
	sqlStr := "select id, file_name, hls_created from video where file_name=?"
	var v Video
	err := db.QueryRow(sqlStr, fileName).Scan(&v.Id, &v.FileName, &v.HlsCreated)
	if err != nil {
		return v, err
	}
	return v, nil
}

func GetVideoFileName(vid ptype.VideoID) (string, error) {
	sqlStr := "select file_name from video where id=?"
	var fileName string
	err := db.QueryRow(sqlStr, vid).Scan(&fileName)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func RemoveVideo(vid ptype.VideoID) error {
	videoPath := setting.GS().Server.HlsPath + fmt.Sprintf("/vid_%d", vid)
	os.RemoveAll(videoPath)
	posterPath := setting.GS().Server.PosterPath + fmt.Sprintf("/vid_%d.jpg", vid)
	os.RemoveAll(posterPath)
	sqlStr := "delete from video where id=?"
	_, err := db.Exec(sqlStr, vid)
	return err
}
