package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/jakewright/home-automation/services/dmx/dmx"
)

var client *dmx.OLAClient

func main() {
	var err error
	client, err = dmx.NewOLAClient("http://ola.local", 9090, 1)
	if err != nil {
		panic(err)
	}

	defer func() {
		update(0, 0, 0, 0)
	}()

	// sig := make(chan os.Signal, 2)
	// signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for i := 0; i < 1; i++ {
		for i := 0; i <= rand.Intn(3); i++ {
			b := rand.Intn(255)
			update(byte(b), 255, 50, 0) // off
			d := rand.Intn(300)
			time.Sleep(time.Millisecond * time.Duration(d))
			update(255, 255, 50, 0) // on
		}

		d := rand.Intn(6000)
		time.Sleep(time.Millisecond * time.Duration(d))
	}
}

func update(brightness, r, g, b byte) {
	if err := client.SetValues(context.TODO(), [512]byte{r, g, b, 0, 0, 0, brightness}); err != nil {
		panic(err)
	}
}
