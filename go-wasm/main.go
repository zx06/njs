//go:build js && wasm

package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"syscall/js"
	"time"
)

func buildinfo() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			slog.Warn("debug.ReadBuildInfo failed")
		}
		slog.Info("buildinfo", slog.Any("info", info))
		return nil
	})
}

func fetchGithub() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		go func() {
			resp, err := http.DefaultClient.Get("https://api.github.com/repos/golang/go/commits?per_page=1")
			if err != nil {
				slog.Error("fetchGithub", slog.Any("err", err))
			}
			defer resp.Body.Close()
			var data any
			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				slog.Error("fetchGithub", slog.Any("err", err))
			}
			slog.Info("fetchGithub", slog.Any("data", data))
		}()
		return nil
	})
}

func main() {
	done := make(chan struct{}, 0)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	slog.Info("Hello Wasm!")
	js.Global().Set("g$buildinfo", buildinfo())
	js.Global().Set("g$fetchGithub", fetchGithub())
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			slog.Info("tick")
		}
	}()
	<-done
}
