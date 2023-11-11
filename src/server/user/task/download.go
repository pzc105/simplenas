package task

import "pnas/prpc"

type OnBtStatus func(prpc.BtStatusRespone)

type Download struct {
	RawTask
	onBtStatus OnBtStatus
}
