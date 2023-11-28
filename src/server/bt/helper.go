package bt

import (
	"context"
	"fmt"
	"pnas/db"
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

func loadInfoHash(id ptype.TorrentID) (*InfoHash, error) {
	sql := `select version, info_hash from torrent where id=?`
	ret := &InfoHash{}
	err := db.QueryRow(sql, id).Scan(&ret.Version, &ret.Hash)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func loadTorrent(btClient *BtClient, id ptype.TorrentID) *Torrent {
	sql := `select id, name, version, info_hash, state, total_size, piece_length, num_pieces, introduce, magnet_uri from torrent where id=?`
	t := &Torrent{}
	err := db.QueryRow(sql, id).Scan(
		&t.base.Id,
		&t.base.Name,
		&t.base.InfoHash.Version,
		&t.base.InfoHash.Hash,
		&t.state,
		&t.base.TotalSize,
		&t.base.PieceLength,
		&t.base.NumPieces,
		&t.base.Introduce,
		&t.base.MagnetUri,
	)
	if err != nil {
		return nil
	}
	t.btClient = btClient
	t.init()
	return t
}

func loadTorrentByInfoHash(btClient *BtClient, infoHash *InfoHash) (*Torrent, error) {
	sql := `select id, name, version, info_hash, state, total_size, piece_length, num_pieces, introduce, magnet_uri
					from torrent where version=? and info_hash=?`
	t := &Torrent{}
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
		&t.base.MagnetUri,
	)
	if err != nil {
		return nil, err
	}
	t.btClient = btClient
	t.init()
	return t, nil
}

func newTorrent(btClient *BtClient, infoHash *InfoHash, magnetUri string) (*Torrent, error) {
	sql := `insert into torrent(version, info_hash, introduce, magnet_uri) values(?, ?, "", ?)`
	r, err := db.Exec(sql, infoHash.Version, infoHash.Hash, magnetUri)
	if err != nil {
		return nil, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	t := &Torrent{
		base: TorrentBase{
			Id:        ptype.TorrentID(id),
			InfoHash:  *infoHash,
			MagnetUri: magnetUri,
		},
	}
	t.btClient = btClient
	t.init()
	return t, nil
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

func delResumeData(infoHash *InfoHash) error {
	_, err := db.GREDIS.HDel(context.Background(), RedisKeyResumeData, getHKey(infoHash)).Result()
	return err
}

func loadResumeData(infoHash *InfoHash) ([]byte, error) {
	d, err := db.GREDIS.HGet(context.Background(), RedisKeyResumeData, getHKey(infoHash)).Result()
	return []byte(d), err
}

func getMagnetByInfoHash(infoHash *InfoHash) (string, error) {
	sql := `select magnet_uri from torrent where info_hash=? and version=?`
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

func IsDownloadAll(st prpc.BtStateEnum) bool {
	return st == prpc.BtStateEnum_seeding
}
