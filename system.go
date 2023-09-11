package main

import (
	"bytes"
	"context"
	"log"
	"strings"
	"time"

	"golang.design/x/clipboard"
)

// Sync piknik and system clipboards on Darwin efficiently.
func SyncClipboards(conf Conf) {

	fromSystem := make(chan string, 1e6)
	fromPiknik := make(chan string, 1e6)

	clipboard.Write(clipboard.FmtText, []byte(""))
	input := bytes.NewReader([]byte(""))

	_, err := RunClient(conf, input, true, false)
	// try again
	for err != nil {
		log.Println(err)
		_, err = RunClient(conf, input, true, false)
	}

	go func() {
		ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
		for data := range ch {
			fromSystem <- string(data)
		}
	}()

	go func() {
		contentOld, err := RunClient(conf, nil, false, false)
		for err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
			contentOld, err = RunClient(conf, nil, false, false)
		}
		messageOld := string(contentOld)
		for {
			time.Sleep(100 * time.Millisecond)
			contentNew, err := RunClient(conf, nil, false, false)
			for err != nil {
				log.Println(err)
				time.Sleep(2 * time.Second)
				contentNew, err = RunClient(conf, nil, false, false)
			}
			messageNew := string(contentNew)

			// send if update
			if messageOld != messageNew {
				messageOld = messageNew
				fromPiknik <- string(messageNew)
			}
		}
	}()

	for {
		select {
		case out := <-fromSystem:
			out = strings.TrimSuffix(out, "\n")
			input := bytes.NewReader([]byte(out))
			_, err = RunClient(conf, input, true, false)

			for err != nil {
				log.Println(err)
				time.Sleep(2 * time.Second)
				_, err = RunClient(conf, input, true, false)
			}

		case out := <-fromPiknik:
			out = strings.TrimSuffix(out, "\n")
			clipboard.Write(clipboard.FmtText, []byte(out))
		}
	}
}
