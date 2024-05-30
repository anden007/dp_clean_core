package part

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/anden007/dp_clean_core/misc"
	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/libkv/store"
	etcdClient "github.com/rpcxio/rpcx-etcd/client"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/share"
	"github.com/spf13/viper"
)

type IRpcxService interface {
	InitConfig() (err error)
	InitRpcx() (err error)
	RegisterName(serviceName string, rcvr interface{}, metadata string) (err error)
	RegisterFunction(servicePath string, fn interface{}, metadata string) (err error)
	RegisterFunctionName(servicePath string, fnName string, fn interface{}, metadata string) (err error)
	GetClient(serviceName string, callback func(clt client.XClient) error) (err error)
	Serve() (err error)
}

type RpcxService struct {
	serviceClientPool *sync.Map
	rpcxServer        *server.Server
	// 注册中心地址
	etcdAddress string
	// 注册中心路径
	basePath string
	// 服务地址，注册中心给客户端用的
	serviceAddress string
	// 实际响应地址，给本机监听客户端请求用的
	listenAddress string
}

func NewRpcxRegistry() (registry IRpcxService) {
	return &RpcxService{
		serviceClientPool: new(sync.Map),
		rpcxServer:        server.NewServer(),
	}
}

func (m *RpcxService) InitConfig() (err error) {
	// loadTime := time.Now()
	// 开起rpcx调试
	if viper.GetBool("rpcx.debug") {
		share.Trace = true
	}
	// if ENV == ENUM_ENV_DEV {
	// 	lib.ServiceLoadInfo("RemoteRPCX", true, loadTime)
	// }
	return
}

func (m *RpcxService) InitRpcx() (err error) {
	port := 0
	// loadTime := time.Now()
	m.etcdAddress = viper.GetString("etcd.address")
	m.basePath = viper.GetString("etcd.base_path")
	listeningIP := os.Getenv("PUBLIC_IP_ADDRESS")
	if listeningIP == "" {
		listeningIP = "127.0.0.1"
	}
	if port == 0 {
		if port, err = misc.GetRndUnUsePortNumber(); err != nil {
			panic(err.Error())
		}
		fmt.Printf("注册服务:%s:%d\n", listeningIP, port)
	}
	m.serviceAddress = fmt.Sprintf("%s:%d", listeningIP, port)
	m.listenAddress = fmt.Sprintf("0.0.0.0:%d", port)
	etcdPlugin := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%s", m.serviceAddress),
		EtcdServers:    []string{m.etcdAddress},
		BasePath:       m.basePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
		Options: &store.Config{
			Username: viper.GetString("etcd.user"),
			Password: viper.GetString("etcd.password"),
		},
	}
	if err = etcdPlugin.Start(); err != nil {
		panic(err.Error())
	}
	m.rpcxServer.Plugins.Add(etcdPlugin)
	// if lib.ENV == lib.ENUM_ENV_DEV {
	// 	lib.ServiceLoadInfo(fmt.Sprintf("RemoteRPCX listening on: %s", m.serviceAddress), true, loadTime)
	// }
	return err
}

func (m *RpcxService) RegisterName(serviceName string, rcvr interface{}, metadata string) (err error) {
	err = m.rpcxServer.RegisterName(serviceName, rcvr, metadata)
	return
}

func (m *RpcxService) RegisterFunction(servicePath string, fn interface{}, metadata string) (err error) {
	err = m.rpcxServer.RegisterFunction(servicePath, fn, metadata)
	return
}

func (m *RpcxService) RegisterFunctionName(servicePath string, fnName string, fn interface{}, metadata string) (err error) {
	err = m.rpcxServer.RegisterFunctionName(servicePath, fnName, fn, metadata)
	return
}

func (m *RpcxService) GetClient(serviceName string, callback func(clt client.XClient) error) (err error) {
	var rpcxClient client.XClient
	if clientItem, exists := m.serviceClientPool.Load(serviceName); exists {
		rpcxClient = clientItem.(client.XClient)
	} else {
		d, _ := etcdClient.NewEtcdV3Discovery(m.basePath, serviceName, []string{m.etcdAddress}, true, &store.Config{
			Username: viper.GetString("etcd.user"),
			Password: viper.GetString("etcd.password"),
		})
		rpcxClient = client.NewXClient(serviceName, client.Failtry, client.RandomSelect, d, client.DefaultOption)
		m.serviceClientPool.Store(serviceName, rpcxClient)
	}
	err = callback(rpcxClient)
	return
}

func (m *RpcxService) Serve() (err error) {
	go func() { // rpxc监听地址
		err = m.rpcxServer.Serve("tcp", m.listenAddress)
	}()
	return
}
