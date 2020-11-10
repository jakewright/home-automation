package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jakewright/home-automation/services/dmx/dmx"
)

var client *dmx.OLAClient

var r, g, b byte

func main() {
	var err error
	client, err = dmx.NewOLAClient("http://ola.local", 9090, 1)
	if err != nil {
		panic(err)
	}

	defer func() {
		update(0)
	}()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Change the colour every 20 minutes
	go func() {
		for {
			r, g, b = 0, 255, 0
			update(255)
			time.Sleep(time.Minute * 5)
			r, g, b = 255, 50, 0
			update(255)
			time.Sleep(time.Minute * 10)
		}
	}()

	for {
		select {
		case <-sig:
			// Exit when a signal is received
			return
		default:
			// Flicker off then on up to 3 times in a row
			// for i := 0; i <= rand.Intn(3); i++ {
			// 	b := rand.Intn(255)
			// 	update(byte(b)) // off
			// 	d := 80 + rand.Intn(200)
			// 	time.Sleep(time.Millisecond * time.Duration(d))
			// 	update(255) // on
			// }

			// Small flickers for a while
			// for i := 0; i < 400+rand.Intn(400); i++ {
			// 	b := 230 + rand.Intn(25)
			// 	update(byte(b))
			// 	d := 50 + rand.Intn(50)
			// 	time.Sleep(time.Millisecond * time.Duration(d))
			// }
		}
	}
}

func update(brightness byte) {
	if err := client.SetValues(context.TODO(), [512]byte{r, g, b, 0, 0, 0, brightness}); err != nil {
		fmt.Printf("error: %v", err)
	}
}
