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

	for _, opt := range opts {
		opt.apply(&bt.opts)
	}

	initClient := func() {
		if len(setting.GS().Bt.BtClientAddress) == 0 {
			log.Error("[bt] empty address")
			return
		}
		if bt.conn!=nil && setting.GS().Bt.BtClientAddress == bt.conn.Target(){
			return
		}
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
		if bt.conn != nil {
			bt.conn.Close()
		}
		bt.conn = conn
		bt.BtServiceClient = prpc.NewBtServiceClient(conn)
		go bt.handleStatus(conn)
		go bt.handleConState(conn)
		go bt.handleFileCompleted(conn)
	}

	initClient()

	setting.AddOnCfgChangeFun("bt_client", initClient)

	bt.closeCtx, bt.closeFunc = context.WithCancel(context.Background())
}

func (bt *BtClient) Close() {
	bt.closeFunc()
	bt.conn.Close()
	bt.wg.Wait()
}

func (bt *BtClient) handleConState(conn *grpc.ClientConn) {
	bt.wg.Add(1)
	defer bt.wg.Done()

	for {
		lst := conn.GetState()
		if !conn.WaitForStateChange(bt.closeCtx, lst) {
			return
		}
		nst := conn.GetState()
		if lst != nst && nst == connectivity.Ready {
			if bt.opts.onConnect != nil {
				bt.opts.onConnect()
			}
		} else if nst == connectivity.TransientFailure {
			log.Warn("[bt] disconnect")
		}
	}
}

func (bt *BtClient) handleStatus(conn *grpc.ClientConn) {
	bt.wg.Add(1)
	defer bt.wg.Done()

	defer func() {
		if conn.GetState() == connectivity.Shutdown {
			return
		}
		time.Sleep(1 * time.Second)
		go bt.handleStatus(conn)
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

func (bt *BtClient) handleFileCompleted(conn *grpc.ClientConn) {
	bt.wg.Add(1)
	defer bt.wg.Done()

	defer func() {
		if conn.GetState() == connectivity.Shutdown {
			return
		}
		time.Sleep(1 * time.Second)
		go bt.handleFileCompleted(conn)
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
