package app

import (
	"common/config"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run 启动程序 启动Grpc服务 启动http服务 启用日志 启用数据库
func Run(ctx context.Context) error {
	// 1.做日志库 info err fatal debug

	// 2. etcd注册中心 把grpc服务注册到etcd中，客户端访问的时候，通过etcd获取grpc地址

	// 启动grpc服务端
	server := grpc.NewServer()
	go func() {
		listen, err := net.Listen("tcp", config.Conf.Grpc.Addr)
		if err != nil {
			log.Fatalf("user grpc server listen err:%v\n", err)
		}
		// 注册grpc service 需要数据库 mongo redis
		// 初始化数据库管理
		// 阻塞操作
		err = server.Serve(listen)
		if err != nil {
			log.Fatalf("user grpc server run failed err:%v\n", err)
		}
	}()
	stop := func() {
		server.Stop()
		// other
		time.Sleep(3 * time.Second)
		fmt.Println("stop app finish")
	}
	// 期望有一个优雅启停 遇到中断信号 退出 终止 挂断
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
	for {
		select {
		case <-ctx.Done():
			stop()
			// time out
			return nil
		case s := <-c:
			switch s {
			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
				stop()
				log.Println("user app quit")
				return nil
			case syscall.SIGHUP:
				stop()
				log.Println("user app SIGHUP")
				return nil
			default:
				return nil
			}
		}
	}
	return nil
}
