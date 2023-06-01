package main

import (
	"context"
	"fmt"
	"github.com/ava-labs/avalanchego/vms/rpcchainvm/grpcutils"
	"net/http"
	"os"
	"time"

	httppb "github.com/ava-labs/avalanchego/proto/pb/http"
	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/vms/rpcchainvm/ghttp"
	abci "github.com/consideritdone/landslidecore/abci/types"
	"github.com/consideritdone/landslidecore/libs/log"
	ctypes "github.com/consideritdone/landslidecore/rpc/core/types"
	"github.com/gorilla/rpc/v2"
)

const (
	Name = "LandslideCore"
)

type VM struct {
	tmLogger log.Logger
}

type LocalService struct {
	vm *VM
}

func (s *LocalService) ABCIInfo(_ *http.Request, _ *struct{}, reply *ctypes.ResultABCIInfo) error {
	reply.Response = abci.ResponseInfo{}
	return nil
}

type Service interface {
	ABCIService
}

type ABCIService interface {
	// Reading from abci app
	ABCIInfo(_ *http.Request, _ *struct{}, reply *ctypes.ResultABCIInfo) error
}

func (vm *VM) CreateHandlers(ctx context.Context) (map[string]*common.HTTPHandler, error) {
	//TODO: why we need mux??
	//mux := http.NewServeMux()
	//rpcLogger := vm.tmLogger.With("module", "rpc-server")
	//rpcserver.RegisterRPCFuncs(mux, rpccore.Routes, rpcLogger)

	server := rpc.NewServer()
	//server.RegisterCodec(json.NewCodec(), "application/json")
	//server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	if err := server.RegisterService(NewService(vm), Name); err != nil {
		return nil, err
	}

	return map[string]*common.HTTPHandler{
		"": {
			//TODO: lock options???
			LockOptions: common.WriteLock,
			Handler:     server,
		},
	}, nil
}

func NewService(vm *VM) Service {
	return &LocalService{vm}
}

func main() {

	vm := VM{
		tmLogger: log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
	}
	handlers, err := vm.CreateHandlers(context.Background())
	if err != nil {
		vm.tmLogger.Error("error during handler creation", err)
	}
	for prefix, h := range handlers {
		handler := h

		//vm.tmLogger.Info("prefix", prefix)

		serverListener, err := grpcutils.NewListener()
		if err != nil {
			vm.tmLogger.Error("error during grpcutils listener creation", prefix, err)
		}
		fmt.Println("Listener", serverListener.Addr())
		server := grpcutils.NewServer()
		httppb.RegisterHTTPServer(server, ghttp.NewServer(handler.Handler))

		// Start HTTP service
		go grpcutils.Serve(serverListener, server)
	}
	for {
		time.Sleep(time.Second)
	}
}
