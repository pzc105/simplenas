package bt

import (
	"pnas/db"
	"pnas/log"
	"pnas/prpc"
)

func TranInfoHash(info *prpc.InfoHash) InfoHash {
	return InfoHash{
		Version: info.GetVersion(),
		Hash:    string(info.GetHash()),
	}
}

func GetInfoHash(infoHash *InfoHash) *prpc.InfoHash {
	return &prpc.InfoHash{
		Version: infoHash.Version,
		Hash:    []byte(infoHash.Hash),
	}
}

func LoadTorrentResumeData(infoHash InfoHash) ([]byte, error) {
	sql := `select resume_data from torrent where info_hash=? and version=?`
	var resumeData []byte
	err := db.QueryRow(sql, infoHash.Hash, infoHash.Version).Scan(&resumeData)
	return resumeData, err
}

func LoadDownloadingTorrent() [][]byte {
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
