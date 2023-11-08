package service

import (
	"fmt"
	"pnas/bt"
	"pnas/prpc"
)

func TranInfoHash(info *prpc.InfoHash) bt.InfoHash {
	return bt.InfoHash{
		Version: info.GetVersion(),
		Hash:    string(info.GetHash()),
	}
}

func GetInfoHash(infoHash *bt.InfoHash) *prpc.InfoHash {
	return &prpc.InfoHash{
		Version: infoHash.Version,
		Hash:    []byte(infoHash.Hash),
	}
}

func getItemRoomKey(room *prpc.Room) string {
	if room.GetType() == prpc.Room_Category {
		return "item_" + fmt.Sprint(room.GetId())
	}
	if room.GetType() == prpc.Room_Danmaku {
		return "danmaku_" + fmt.Sprint(room.GetId())
	}
	return ""
}
