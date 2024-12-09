package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	// "os"

	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context) error{
	s := &http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	//別ゴルーチンでHTTPサーバーを起動する。
	eg.Go(func()error{
		//http.ErrServerClosedは
		//http.Server.Shutdown()が正常終了したことを示すので異常ではない。
		if err := s.ListenAndServe(); err != nil &&
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