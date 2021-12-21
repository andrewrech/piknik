package main

import (
	"bytes"
	"context"
	"time"

	"golang.design/x/clipboard"
)

func systemClipboardChan() (out chan string) {

	go func() {
		out = make(chan string, 1e6)
		ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
		for data := range ch {
			out <- string(data)
		}
	}()

	return out
}

func piknikChan(conf Conf) (out chan string) {

	out = make(chan string, 1e6)
	go func() {
		contentOld := string(RunClient(conf, nil, false, false))
		for {
			time.Sleep(100 * time.Millisecond)
			contentNew := string(RunClient(conf, nil, false, false))

			// send if update
			if contentOld != contentNew {
				contentOld = contentNew
				out <- string(contentNew)
			}
		}
	}()

	return out
}

// Sync piknik and system clipboards on Darwin efficiently.
func SyncClipboards(conf Conf, fromSystem chan string, fromPiknik chan string) {

	clipboard.Write(clipboard.FmtText, []byte(""))
	_ = RunClient(conf, bytes.NewReader([]byte("")), true, false)

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
