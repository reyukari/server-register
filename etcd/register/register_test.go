package register

import (
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestNewServiceRegister(t *testing.T) {
	s, err := NewRegister(
		SetName("node.srv.app"),
		SetAddress("127.0.0.1:123123"),
		SetWeight("1"),
		SetEtcdConf(clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: time.Second * 5,
		}),
	)
	if err != nil {
		panic(err)
	}
	c := make(chan os.Signal, 1)
	go func() {
		if s.ListenKeepAliveChan() {
			c <- syscall.SIGQUIT
		}
	}()
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for a := range c {
		switch a {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			zap.S().Info("退出")
			_ = s.Close()
			return
		default:
			return
		}
	}
}
