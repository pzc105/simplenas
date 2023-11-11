package bt

import (
	"context"
	"io"
	"pnas/log"
	"pnas/prpc"
	"pnas/setting"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type BtClient struct {
	prpc.BtServiceClient
	conn *grpc.ClientConn
	opts btClientOpts

	closeCtx  context.Context
	closeFunc context.CancelFunc
	wg        sync.WaitGroup
}

type btClientOpts struct {
	onStatus        func(*prpc.BtStatusRespone)
	onFileCompleted func(*prpc.FileCompletedRes)
	onConnect       func()
}

type BtClientOpt interface {
	apply(*btClientOpts)
}

type funcBtClientOpt struct {
	do func(opts *btClientOpts)
}

func (f *funcBtClientOpt) apply(opts *btClientOpts) {
	f.do(opts)
}

func WithOnStatus(onStatus func(*prpc.BtStatusRespone)) *funcBtClientOpt {
	return &funcBtClientOpt{
		do: func(opts *btClientOpts) {
			opts.onStatus = onStatus
		},
	}
}

func WithOnConnect(onConnect func()) *funcBtClientOpt {
	return &funcBtClientOpt{
		do: func(opts *btClientOpts) {
			opts.onConnect = onConnect
		},
	}
}

func WithOnFileCompleted(onFileCompleted func(*prpc.FileCompletedRes)) *funcBtClientOpt {
	return &funcBtClientOpt{
		do: func(opts *btClientOpts) {
			opts.onFileCompleted = onFileCompleted
		},
	}
}

func (bt *BtClient) Init(opts ...BtClientOpt) {
	bc := backoff.DefaultConfig
	bc.MaxDelay = time.Second * 3
	conn, err := grpc.Dial(setting.GS().Bt.BtClientAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           bc,
			MinConnectTimeout: time.Second * 5,
		}))
	if err != nil {
		log.Error("failed to connect bt")
		return
	}
	for _, opt := range opts {
		opt.apply(&bt.opts)
	}

	bt.conn = conn
	bt.BtServiceClient = prpc.NewBtServiceClient(conn)
	bt.closeCtx, bt.closeFunc = context.WithCancel(context.Background())
	go bt.handleStatus()
	go bt.handleConState()
	go bt.handleFileCompleted()
}

func (bt *BtClient) Close() {
	bt.closeFunc()
	bt.conn.Close()
	bt.wg.Wait()
}

func (bt *BtClient) handleConState() {
	bt.wg.Add(1)
	defer bt.wg.Done()

	for {
		lst := bt.conn.GetState()
		if !bt.conn.WaitForStateChange(bt.closeCtx, lst) {
			return
		}
		nst := bt.conn.GetState()
		if lst != nst && nst == connectivity.Ready {
			if bt.opts.onConnect != nil {
				bt.opts.onConnect()
			}
		} else if nst == connectivity.TransientFailure {
			log.Warn("[bt] disconnect")
		}
	}
}

func (bt *BtClient) handleStatus() {
	bt.wg.Add(1)
	defer bt.wg.Done()

	defer func() {
		if bt.conn.GetState() == connectivity.Shutdown {
			return
		}
		time.Sleep(1 * time.Second)
		go bt.handleStatus()
	}()

	stream, err := bt.OnBtStatus(bt.closeCtx)
	if err != nil {
		return
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			stream.CloseSend()
			break
		}
		if err != nil {
			stream.CloseSend()
			break
		}

		if bt.opts.onStatus != nil {
			bt.opts.onStatus(response)
		}
	}
}

func (bt *BtClient) handleFileCompleted() {
	bt.wg.Add(1)
	defer bt.wg.Done()

	defer func() {
		if bt.conn.GetState() == connectivity.Shutdown {
			return
		}
		time.Sleep(1 * time.Second)
		go bt.handleFileCompleted()
	}()

	stream, err := bt.OnFileCompleted(bt.closeCtx)
	if err != nil {
		return
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			stream.CloseSend()
			break
		}
		if err != nil {
			stream.CloseSend()
			break
		}

		if bt.opts.onFileCompleted != nil {
			bt.opts.onFileCompleted(response)
		}
	}
}
