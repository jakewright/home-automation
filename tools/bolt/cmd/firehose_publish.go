package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/spf13/cobra"

	"github.com/jakewright/home-automation/libraries/go/firehose"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

var (
	firehosePublishCmd = &cobra.Command{
		Use:   "publish [channel] [json payload]",
		Short: "publish raw JSON messages to the Firehose",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				_ = cmd.Usage()
				return
			}

			fmt.Println(args[1])

			var msg interface{}
			if err := json.Unmarshal([]byte(args[1]), &msg); err != nil {
				output.Fatal("Failed to unmarshal JSON: %v", err)
			}

			addr := "localhost:6379"
			redisClient := redis.NewClient(&redis.Options{
				Addr:            addr,
				Password:        "",
				DB:              0,
				MaxRetries:      1,
				MinRetryBackoff: time.Second,
				MaxRetryBackoff: time.Second * 5,
			})

			defer func() { _ = redisClient.Close() }()

			c := firehose.NewStreamsClient(redisClient)

			if err := c.Publish(args[0], msg); err != nil {
				output.Fatal("Failed to publish: %v", err)
			}

			output.Info("Messaged published to channel %q", args[0]).Success()
		},
	}
)

func init() {
	firehoseCmd.AddCommand(firehosePublishCmd)
}
