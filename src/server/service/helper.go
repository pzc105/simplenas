package service

import (
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
