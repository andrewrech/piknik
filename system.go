package main

import (
	"bytes"
	"context"
	"time"

	"golang.design/x/clipboard"
)

// Sync piknik and system clipboards on Darwin efficiently.
func SyncClipboards(conf Conf) {

	fromSystem := make(chan string, 1e6)
	fromPiknik := make(chan string, 1e6)

	clipboard.Write(clipboard.FmtText, []byte(""))
	input := bytes.NewReader([]byte(""))
	_ = RunClient(conf, input, true, false)

	go func() {
		ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
		for data := range ch {
			fromSystem <- string(data)
		}
	}()

	go func() {
		contentOld := string(RunClient(conf, nil, false, false))
		for {
			time.Sleep(100 * time.Millisecond)
			contentNew := string(RunClient(conf, nil, false, false))

			// send if update
			if contentOld != contentNew {
				contentOld = contentNew
				fromPiknik <- string(contentNew)
			}
		}
	}()

	for {
		select {
		case out := <-fromSystem:
			input := bytes.NewReader([]byte(out))
			_ = RunClient(conf, input, true, false)

		case out := <-fromPiknik:
			clipboard.Write(clipboard.FmtText, []byte(out))
		}
	}
}
