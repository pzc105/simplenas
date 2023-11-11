package bt

import (
	"context"
	"fmt"
	"pnas/db"
	"pnas/log"
	"pnas/prpc"
	"pnas/ptype"
	"strings"
)

func TranInfoHash(info *prpc.InfoHash) *InfoHash {
	return &InfoHash{
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

func loadTorrent(btClient *BtClient, id ptype.TorrentID) *Torrent {
	sql := `select name, version, info_hash, state, total_size, piece_length, num_pieces, introduce from torrent where id=?`
	var t Torrent
	t.base.Id = id
	err := db.QueryRow(sql, id).Scan(
		&t.base.Name,
		&t.base.InfoHash.Version,
		&t.base.InfoHash.Hash,
		&t.state,
		&t.base.TotalSize,
		&t.base.PieceLength,
		&t.base.NumPieces,
		&t.base.Introduce,
	)
	if err != nil {
		return nil
	}
	t.btClient = btClient
	t.init()
	return &t
}

func loadTorrentByInfoHash(btClient *BtClient, infoHash *InfoHash) *Torrent {
	sql := `select id, name, version, info_hash, state, total_size, piece_length, num_pieces, introduce 
					from torrent where version=? and info_hash=?`
	var t Torrent
	err := db.QueryRow(sql, infoHash.Version, infoHash.Hash).Scan(
		&t.base.Id,
		&t.base.Name,
		&t.base.InfoHash.Version,
		&t.base.InfoHash.Hash,
		&t.state,
		&t.base.TotalSize,
		&t.base.PieceLength,
		&t.base.NumPieces,
		&t.base.Introduce,
	)
	if err != nil {
		return nil
	}
	t.btClient = btClient
	t.init()
	return &t
}

func newTorrent(btClient *BtClient, infoHash *InfoHash) *Torrent {
	sql := `insert into torrent(version, info_hash, introduce) values(?, ?, "")`
	r, err := db.Exec(sql, infoHash.Version, infoHash.Hash)
	if err != nil {
		return nil
	}
	id, err := r.LastInsertId()
	if err != nil {
		return nil
	}
	t := &Torrent{
		base: TorrentBase{
			Id:       ptype.TorrentID(id),
			InfoHash: *infoHash,
		},
	}
	t.btClient = btClient
	t.init()
	return t
}

func deleteUserTorrentids(tid ptype.TorrentID, ids ...ptype.UserID) error {
	if len(ids) == 0 {
		return nil
	}
	var cond strings.Builder
	for _, uid := range ids {
		if cond.Len() == 0 {
			cond.WriteString(fmt.Sprint(uid))
		} else {
			cond.WriteString("," + fmt.Sprint(uid))
		}
	}
	sql := fmt.Sprintf(`delete from user_torrent where torrent_id=? and user_id in (%s)`, cond.String())
	_, err := db.Exec(sql, tid)
	return err
}

const (
	RedisKeyResumeData      = "torrent_resume"
	RedisKeyBtSessionParams = "bt_session"
)

func getHKey(infoHash *InfoHash) string {
	return fmt.Sprint(infoHash.Version) + infoHash.Hash
}

func saveResumeData(infoHash *InfoHash, data []byte) error {
	_, err := db.GREDIS.HSet(context.Background(), RedisKeyResumeData, getHKey(infoHash), data).Result()
	return err
}

func loadResumeData(infoHash *InfoHash) ([]byte, error) {
	d, err := db.GREDIS.HGet(context.Background(), RedisKeyResumeData, getHKey(infoHash)).Result()
	return []byte(d), err
}

func getMagnetByInfoHash(infoHash *InfoHash) (string, error) {
	sql := `select magnet_uri from magnet where info_hash=? and version=?`
	var ret string
	err := db.QueryRow(sql, infoHash.Hash, infoHash.Version).Scan(&ret)
	return ret, err
}

func saveBtSessionParams(data []byte) error {
	_, err := db.GREDIS.Set(context.Background(), RedisKeyBtSessionParams, data, 0).Result()
	return err
}

func loadBtSessionParams() ([]byte, error) {
	d, err := db.GREDIS.Get(context.Background(), RedisKeyBtSessionParams).Result()
	return []byte(d), err
}

func saveMagnetUri(infoHash *InfoHash, uri string) {
	sql := "insert into magnet(version, info_hash, magnet_uri) values(?, ?, ?) on duplicate key update magnet_uri=values(magnet_uri)"
	_, err := db.Exec(sql, infoHash.Version, infoHash.Hash, uri)
	if err != nil {
		log.Debugf("failed to save err: %v", err)
	}
}
