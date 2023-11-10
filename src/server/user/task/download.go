package task

import "pnas/prpc"

type OnBtStatus func(prpc.StatusRespone)

type Download struct {
	RawTask
	onBtStatus OnBtStatus
}
