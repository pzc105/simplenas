package service

import (
	"errors"
	"net/url"
	"pnas/bt"
	"pnas/category"
	"pnas/prpc"
	"pnas/user"
	"strconv"
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

func IsShared(shares SharesInterface, um *user.UserManger, queryParams url.Values) bool {
	itemIdTmp, _ := strconv.ParseInt(queryParams.Get("itemid"), 10, 64)
	itemId := category.ID(itemIdTmp)
	shareid := queryParams.Get("shareid")
	si, err := shares.GetShareItemInfo(shareid)
	if err != nil {
		return false
	}
	if !um.IsItemShared(si.ShareItemInfo.ItemId, category.ID(itemId)) {
		return false
	}
	return true
}

func GetSharedItemInfo(shares SharesInterface, um *user.UserManger, queryParams url.Values) (user.ID, category.ID, error) {
	itemIdTmp, _ := strconv.ParseInt(queryParams.Get("itemid"), 10, 64)
	itemId := category.ID(itemIdTmp)
	shareid := queryParams.Get("shareid")
	si, err := shares.GetShareItemInfo(shareid)
	if err != nil {
		return -1, -1, errors.New("not found shared item info")
	}
	if !um.IsItemShared(si.ShareItemInfo.ItemId, category.ID(itemId)) {
		return -1, -1, errors.New("item is not shared")
	}
	return si.UserId, itemId, nil
}
