package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
    "os"
	"golang.org/x/sync/errgroup"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v",p,err)
	}
	//エラー処理
	if err := run(context.Background(),l); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, l net.Listener) error{
	//サーバー作成
	s := &http.Server{
		//引数で受け取ったnet.Listenerを利用するのでAddrフィールドは設定しない
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	//別ゴルーチン（並行処理？）でHTTPサーバーを起動する。
	eg.Go(func()error{
		//http.ErrServerClosedは
		//http.Server.Shutdown()が正常終了したことを示すので異常ではない。
		if err := s.Serve(l); err != nil &&
		err != http.ErrServerClosed{
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})
	//チャンネルからの通知（終了通知）を待機する
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	return eg.Wait()
}