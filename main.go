package main

import (
	"context"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/ips"
	abci "github.com/consideritdone/landslidecore/abci/types"
	"github.com/consideritdone/landslidecore/libs/log"
	ctypes "github.com/consideritdone/landslidecore/rpc/core/types"
	rpctypes "github.com/consideritdone/landslidecore/rpc/jsonrpc/types"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"path"
	"time"
	"vm_rpc_server/rpc/jsonrpc/server"
	"vm_rpc_server/utils"
)

//import (
//	"context"
//	"fmt"
//	"github.com/ava-labs/avalanchego/vms/rpcchainvm/grpcutils"
//	"net/http"
//	"os"
//	"time"
//
//	httppb "github.com/ava-labs/avalanchego/proto/pb/http"
//	"github.com/ava-labs/avalanchego/snow/engine/common"
//	"github.com/ava-labs/avalanchego/vms/rpcchainvm/ghttp"
//	abci "github.com/consideritdone/landslidecore/abci/types"
//	"github.com/consideritdone/landslidecore/libs/log"
//	ctypes "github.com/consideritdone/landslidecore/rpc/core/types"
//	"github.com/gorilla/rpc/v2"
//)
//
//const (
//	Name = "LandslideCore"
//)
//
//type VM struct {
//	tmLogger log.Logger
//}
//
//type LocalService struct {
//	vm *VM
//}
//
//func (s *LocalService) ABCIInfo(_ *http.Request, _ *struct{}, reply *ctypes.ResultABCIInfo) error {
//	reply.Response = abci.ResponseInfo{}
//	return nil
//}
//
//type Service interface {
//	ABCIService
//}
//
//type ABCIService interface {
//	// Reading from abci app
//	ABCIInfo(_ *http.Request, _ *struct{}, reply *ctypes.ResultABCIInfo) error
//}
//
//func (vm *VM) CreateHandlers(ctx context.Context) (map[string]*common.HTTPHandler, error) {
//	//TODO: why we need mux??
//	//mux := http.NewServeMux()
//	//rpcLogger := vm.tmLogger.With("module", "rpc-server")
//	//rpcserver.RegisterRPCFuncs(mux, rpccore.Routes, rpcLogger)
//
//	server := rpc.NewServer()
//	//server.RegisterCodec(json.NewCodec(), "application/json")
//	//server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
//	if err := server.RegisterService(NewService(vm), Name); err != nil {
//		return nil, err
//	}
//
//	return map[string]*common.HTTPHandler{
//		"": {
//			//TODO: lock options???
//			LockOptions: common.WriteLock,
//			Handler:     server,
//		},
//	}, nil
//}
//
//func NewService(vm *VM) Service {
//	return &LocalService{vm}
//}
//
//func main() {
//
//	vm := VM{
//		tmLogger: log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
//	}
//	handlers, err := vm.CreateHandlers(context.Background())
//	if err != nil {
//		vm.tmLogger.Error("error during handler creation", err)
//	}
//	for prefix, h := range handlers {
//		handler := h
//
//		//vm.tmLogger.Info("prefix", prefix)
//
//		serverListener, err := grpcutils.NewListener()
//		if err != nil {
//			vm.tmLogger.Error("error during grpcutils listener creation", prefix, err)
//		}
//		fmt.Println("Listener", serverListener.Addr())
//		server := grpcutils.NewServer()
//		httppb.RegisterHTTPServer(server, ghttp.NewServer(handler.Handler))
//
//		// Start HTTP service
//		go grpcutils.Serve(serverListener, server)
//	}
//	for {
//		time.Sleep(time.Second)
//	}
//}

const (
	Name = "LandslideCore"
)

const baseURL = "/ext"

type VM struct {
	tmLogger log.Logger
}

type LocalService struct {
	vm *VM
}

func NewService(vm *VM) Service {
	return &LocalService{vm}
}

// Routes is a map of available routes.
func (vm *VM) RPCRoutes() map[string]*server.RPCFunc {
	vmTMService := NewService(vm)
	//defaultEndpoint := path.Join(constants.ChainAliasPrefix, "chain-id")
	//url := fmt.Sprintf("%s/%s", baseURL, defaultEndpoint)
	return map[string]*server.RPCFunc{
		//// subscribe/unsubscribe are reserved for websocket events.
		//"subscribe":       rpc.NewWSRPCFunc(Subscribe, "query"),
		//"unsubscribe":     rpc.NewWSRPCFunc(Unsubscribe, "query"),
		//"unsubscribe_all": rpc.NewWSRPCFunc(UnsubscribeAll, ""),
		//
		//// info API
		//"health":               rpc.NewRPCFunc(Health, ""),
		//"status":               rpc.NewRPCFunc(Status, ""),
		//"net_info":             rpc.NewRPCFunc(NetInfo, ""),
		//"blockchain":           rpc.NewRPCFunc(BlockchainInfo, "minHeight,maxHeight"),
		//"genesis":              rpc.NewRPCFunc(Genesis, ""),
		//"genesis_chunked":      rpc.NewRPCFunc(GenesisChunked, "chunk"),
		//"block":                rpc.NewRPCFunc(Block, "height"),
		//"block_by_hash":        rpc.NewRPCFunc(BlockByHash, "hash"),
		//"block_results":        rpc.NewRPCFunc(BlockResults, "height"),
		//"commit":               rpc.NewRPCFunc(Commit, "height"),
		//"check_tx":             rpc.NewRPCFunc(CheckTx, "tx"),
		//"tx":                   rpc.NewRPCFunc(Tx, "hash,prove"),
		//"tx_search":            rpc.NewRPCFunc(TxSearch, "query,prove,page,per_page,order_by"),
		//"block_search":         rpc.NewRPCFunc(BlockSearch, "query,page,per_page,order_by"),
		//"validators":           rpc.NewRPCFunc(Validators, "height,page,per_page"),
		//"dump_consensus_state": rpc.NewRPCFunc(DumpConsensusState, ""),
		//"consensus_state":      rpc.NewRPCFunc(ConsensusState, ""),
		//"consensus_params":     rpc.NewRPCFunc(ConsensusParams, "height"),
		//"unconfirmed_txs":      rpc.NewRPCFunc(UnconfirmedTxs, "limit"),
		//"num_unconfirmed_txs":  rpc.NewRPCFunc(NumUnconfirmedTxs, ""),
		//
		//// tx broadcast API
		//"broadcast_tx_commit": rpc.NewRPCFunc(BroadcastTxCommit, "tx"),
		//"broadcast_tx_sync":   rpc.NewRPCFunc(BroadcastTxSync, "tx"),
		//"broadcast_tx_async":  rpc.NewRPCFunc(BroadcastTxAsync, "tx"),
		//
		//// abci API
		//"abci_query": rpc.NewRPCFunc(ABCIQuery, "path,data,height,prove"),
		//fmt.Sprintf("%s/%s", url, "abci_info"): rpcserver.NewRPCFunc(vmTMService.ABCIInfo, ""),
		"abci_info": server.NewRPCFunc(vmTMService.ABCIInfo, ""),
		//
		//// evidence API
		//"broadcast_evidence": rpc.NewRPCFunc(BroadcastEvidence, "evidence"),
	}
}

//func (s *LocalService) ABCIInfo(_ *http.Request, _ *struct{}, reply *ctypes.ResultABCIInfo) error {
//	reply.Response = abci.ResponseInfo{}
//	return nil
//}

// ABCIInfo gets some info about the application.
// More: https://docs.tendermint.com/master/rpc/#/ABCI/abci_info
func (s *LocalService) ABCIInfo(ctx *rpctypes.Context) (*ctypes.ResultABCIInfo, error) {

	return &ctypes.ResultABCIInfo{Response: abci.ResponseInfo{}}, nil
}

type Service interface {
	ABCIService
}

type ABCIService interface {
	// Reading from abci app
	//TODO
	//ABCIInfo(_ *http.Request, _ *struct{}, reply *ctypes.ResultABCIInfo) error
	ABCIInfo(ctx *rpctypes.Context) (*ctypes.ResultABCIInfo, error)
}

type timeHandler struct {
	format string
}

func (th timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.format)
	w.Write([]byte("The time is: " + tm))
}

func (vm *VM) CreateHandlers(_ context.Context) (map[string]*common.HTTPHandler, error) {
	//server := rpc.NewServer()
	//tmService := NewService(vm)
	//err := rpc.Register(tmService)
	//if err != nil {
	//	return nil, fmt.Errorf("Failed to create vm handlers: failed to register tendermint service: %w ", err)
	//}
	//rpc.HandleHTTP()

	mux := http.NewServeMux()

	// 1) Register regular routes.
	routes := vm.RPCRoutes()
	//rpcserver.RegisterRPCFuncs(mux, routes, vm.tmLogger)

	//USE code below instead
	//for funcName, rpcFunc := range routes {
	//	mux.HandleFunc(fmt.Sprintf("%s/%s", "/ext/bc/chain-id/rpc", funcName), server.MakeHTTPHandler(rpcFunc, vm.tmLogger))
	//}

	// JSONRPC endpoints
	//mux.HandleFunc("/ext/bc/chain-id/rpc", server.HandleInvalidJSONRPCPaths(server.MakeJSONRPCHandler(routes, vm.tmLogger)))//TODO
	mux.HandleFunc("/ext/bc/chain-id/rpc", server.MakeJSONRPCHandler(routes, vm.tmLogger))

	//th := timeHandler{format: time.DateTime}

	return map[string]*common.HTTPHandler{
		"/rpc": {
			LockOptions: common.WriteLock,
			Handler:     mux, //TODO
			//Handler: th,
		},
	}, nil
}

func main() {
	vm := VM{
		tmLogger: log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
	}
	handlers, err := vm.CreateHandlers(context.Background())
	if err != nil {
		vm.tmLogger.Error("error during handler creation", err)
	}
	listenHost := "localhost"
	listenPort := 8080

	//listener, err := rpcserver.Listen("localhost:8080", nil)
	//if err != nil {
	//	vm.tmLogger.Error("error during handler creation", err)
	//}

	//server := grpcutils.NewServer()
	router := utils.NewRouter()
	defaultEndpoint := path.Join(constants.ChainAliasPrefix, "chain-id")
	for endpoint, handler := range handlers {
		//httppb.RegisterHTTPServer(server, ghttp.NewServer(handler.Handler))
		//// Start HTTP service
		//go grpcutils.Serve(serverListener, server)
		url := fmt.Sprintf("%s/%s", baseURL, defaultEndpoint)
		vm.tmLogger.Info("adding route",
			zap.String("url", url),
			zap.String("endpoint", endpoint),
		)
		router.AddRouter(url, endpoint, handler.Handler)
	}

	listenAddress := fmt.Sprintf("%s:%d", listenHost, listenPort)
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		vm.tmLogger.Error("Failed to create listener: ", err)
	}

	ipPort, err := ips.ToIPPort(listener.Addr().String())
	if err != nil {
		vm.tmLogger.Info("HTTP API server listening",
			zap.String("address", listenAddress),
		)
	} else {
		vm.tmLogger.Info("HTTP API server listening",
			zap.String("host", listenHost),
			zap.Uint16("port", ipPort.Port),
		)
	}

	corsHandler := cors.New(cors.Options{
		AllowCredentials: true,
	}).Handler(router)
	gzipHandler := gziphandler.GzipHandler(corsHandler)
	var handler http.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Attach this node's ID as a header
			w.Header().Set("node-id", "id123")
			gzipHandler.ServeHTTP(w, r)
		},
	)

	srv := &http.Server{
		Handler: handler,
	}

	if err := srv.Serve(listener); err != nil {
		vm.tmLogger.Error("Failed to serve listener", err)
	}

	//// Start HTTP service
	//go grpcutils.Serve(listener, server)
}
