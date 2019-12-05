package config

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"

	"github.com/pganalyze/collector/util"
)

type googleLogResource struct {
	ResourceType string            `json:"type"`
	Labels       map[string]string `json:"labels"`
}

type googleLogMessage struct {
	InsertID         string            `json:"insertId"`
	LogName          string            `json:"logName"`
	ReceiveTimestamp string            `json:"receiveTimestamp"`
	Resource         googleLogResource `json:"resource"`
	Severity         string            `json:"severity"`
	TextPayload      string            `json:"textPayload"`
	Timestamp        string            `json:"timestamp"`
}

func (s *Config) setupCloudSQLLogSubscriber() (logger *util.Logger, config ServerConfig) {
	projectID := config.GcpProjectID
	subID := config.GcpPubsubSubscription

	ctx := context.Background()
	var opts []option.ClientOption
	if config.GcpCredentialsFile != "" {
		logger.PrintVerbose("Using GCP credentials file located at: %s", config.GcpCredentialsFile)
		opts = append(opts, option.WithCredentialsFile(config.GcpCredentialsFile))
	}
	client, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		logger.PrintError("Failed to create Google PubSub client: %v", err)
		return
	}

	sub := client.Subscription(subID)
	cctx, _ := context.WithCancel(ctx) // cctx, cancel =
	err = sub.Receive(cctx, func(ctx context.Context, pubsubMsg *pubsub.Message) {
		var msg googleLogMessage
		json.Unmarshal(pubsubMsg.Data, &msg)
		//fmt.Printf("Got message: %+v\n", msg)
		t, _ := time.Parse(time.RFC3339Nano, msg.Timestamp)

		pubsubMsg.Ack()
		s.GcpLogStream <- GcpLogStreamItem{Content: msg.TextPayload, OccurredAt: t}
		// cancel() // Do we need to cancel sometime? (e.g. when main program shuts down?)
	})
	if err != nil {
		logger.PrintError("Failed to receive from Google PubSub: %v", err)
	}
	return
}

func handleCloudSQL() (conf Config) {
	conf.GcpLogStream = make(chan GcpLogStreamItem, streamBufferLen)

	// Iterate over all servers and call log setup handler once per unique subscription name

	return
}
