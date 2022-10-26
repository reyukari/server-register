package etcd_grpc

import (
	"context"
	"fmt"
	"github.com/reyukari/server-register/etcd/etcd-grpc/api"
	"github.com/reyukari/server-register/loadbalence"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	r, err := loadbalence.NewUsageLB(
		loadbalence.SetName(ServerName),
		loadbalence.SetLoadBalancingPolicy(loadbalence.UsageLB),
		loadbalence.SetEtcdConf(clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: time.Second * 5,
		}))
	if err != nil {
		panic(err)
	}
	resolver.Register(r)
	// 连接服务器
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", r.Scheme(), ""),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, loadbalence.UsageLB)),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()
	apiClient := api.NewApiClient(conn)

	res, err := apiClient.ApiTest(context.Background(), &api.Request{Input: "usageLB"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Output)
}
