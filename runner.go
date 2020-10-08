package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	/*
		#cgo darwin LDFLAGS: -framework CoreGraphics
		#if defined(__APPLE__)
		#include <CoreGraphics/CGDisplayConfiguration.h>
		int display_width() {
		return CGDisplayPixelsWide(CGMainDisplayID());
		}
		int display_height() {
		return CGDisplayPixelsHigh(CGMainDisplayID());
		}
		#else
		int display_width() {
		return 1024;
		}
		int display_height() {
		return 768;
		}
		#endif
	*/
	"C"

	"github.com/phayes/freeport"
	"github.com/zserge/lorca"
)

func main() {
	if lorca.LocateChrome() == "" {
		lorca.PromptDownload()
		return
	}

	w, err := lorca.New(
		"data:text/html,<title>Loading...</title>",
		"",
		int(C.display_width()),
		int(C.display_height()),
	)
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" || port == "0" {
		p, err := freeport.GetFreePort()
		if err != nil {
			log.Fatal(err)
		}

		port = strconv.Itoa(p)
		os.Setenv("PORT", port)
	}

	url := strings.Join([]string{"http://localhost:", port}, "")

	srv := os.Getenv("SERVER")
	if srv == "" {
		srv = "makehtml.jar"
	}

	cmd := exec.Command("java", "-jar", srv)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	w.SetBounds(lorca.Bounds{
		WindowState: lorca.WindowStateMaximized,
	})

	defer cmd.Process.Kill()
	defer w.Close()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			_, err := http.Head(url)
			if err == nil {
				break
			}
		}

		w.Load(url)
	}()

	<-w.Done()
}
