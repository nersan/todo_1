package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	// "time"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	// コンテキストとキャンセル関数を用意
	ctx, cancel := context.WithCancel(context.Background())

	// サーバーを別ゴルーチンで起動
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})

	// サーバーにリクエストを送信
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()

	// レスポンスボディを読み取る
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %+v", err)
	}

	// レスポンス内容を検証
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	// サーバーを終了
	cancel()

	// サーバーの終了を待つ
	if err := eg.Wait(); err != nil && err != context.Canceled {
		t.Fatalf("unexpected error: %+v", err)
	}
}
