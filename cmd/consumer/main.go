package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	client, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
		kgo.ConsumerGroup("orders-consumers"),
		kgo.ConsumeTopics("orders"),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		kgo.OnPartitionsAssigned(func(_ context.Context, _ *kgo.Client, assigned map[string][]int32) {
			fmt.Printf("assigned: %v\n", assigned)
		}),
		kgo.OnPartitionsRevoked(func(_ context.Context, _ *kgo.Client, revoked map[string][]int32) {
			fmt.Printf("revoked: %v\n", revoked)
		}),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	fmt.Println("waiting for messages, press Ctrl+C to exit")

	for {
		fetches := client.PollFetches(ctx)
		if ctx.Err() != nil || fetches.IsClientClosed() {
			break
		}

		fetches.EachError(func(topic string, partition int32, err error) {
			log.Printf("fetch error: topic=%s partition=%d err=%v", topic, partition, err)
		})

		fetches.EachRecord(func(r *kgo.Record) {
			fmt.Printf("received: partition=%d offset=%d key=%s value=%s\n",
				r.Partition, r.Offset, r.Key, r.Value)
		})
	}

	fmt.Println("shutting down")
}
