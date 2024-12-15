package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	// コンテキストとキャンセル関数を用意
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// サーバーを別ゴルーチンで起動
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})

	// サーバー起動を待つ
	time.Sleep(100 * time.Millisecond) // 起動を待つための簡易的な遅延（本番環境では適切な同期を推奨）

	// サーバーにリクエストを送信
	in := "message"
	resp, err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Fatalf("failed to get: %+v", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを読み取る
	got, err := io.ReadAll(resp.Body)
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
