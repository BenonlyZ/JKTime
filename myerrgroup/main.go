package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

//传入一个http server 和 struct{}类型的channel便于处理关闭请求
func MyHttpServer(s *http.Server, shutdown chan struct{}) error {
	//注册一个打印helloworld的路由
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	})
	//注册一个处理关闭http server的路由
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		//请求关闭server时,往通道shutdown传入struct{}类型对象值
		shutdown <- struct{}{}
	})
	fmt.Println("http server监听端口", s.Addr)
	return s.ListenAndServe()

}

func main() {
	//初始化一个带取消的Group和context
	g, errctx := errgroup.WithContext(context.Background())
	//初始化一个http server
	server := http.Server{
		Addr: ":8087",
	}
	//初始化通道shutdown
	shutdown := make(chan struct{})

	//起一个goroutine来启动http server
	g.Go(func() error {
		return MyHttpServer(&server, shutdown)
	})

	//起一个goroutine来实现http server的关闭
	g.Go(func() error {
		select {
		//当接收一个signal信号退出程序时(ctrl+C),errctx.Done()通道关闭,则非阻塞,返回该类型零值
		case <-errctx.Done():
			fmt.Printf("http server errgroup exit by: %+v\n", errctx.Err())
		//当客户端请求关闭时,<-shutdown非阻塞
		case <-shutdown:
			fmt.Println("server shutting ...")
		}
		//errctx->shutCtx context, 并返回一个取消方法
		shutCtx, cancel := context.WithCancel(errctx)
		defer cancel()
		//关闭http server
		return server.Shutdown(shutCtx)
	})

	//起一个goroutine来实现SIGINT(ctrl+C)信号的注册和处理
	g.Go(func() error {
		quit := make(chan os.Signal, 0)
		signal.Notify(quit, syscall.SIGINT)

		select {
		//当http server关闭时,调用cancel()取消context,此时errctx.Done()关闭且非阻塞
		case <-errctx.Done():
			return errctx.Err()
		//当接收一个SIGINT型号时,<-quit非阻塞
		case <-quit:
			return errors.Errorf("get os signal: %v", <-quit)

		}
	})

	fmt.Printf("errgroup exiting: %+v\n", g.Wait())
}
