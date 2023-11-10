package utils

import (
	"fmt"
	"pnas/prpc"
)

func GetItemRoomKey(room *prpc.Room) string {
	if room.GetType() == prpc.Room_Category {
		return "item_" + fmt.Sprint(room.GetId())
	}
	if room.GetType() == prpc.Room_Danmaku {
		return "danmaku_" + fmt.Sprint(room.GetId())
	}
	return ""
}
