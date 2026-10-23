package main

import (
	"context"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	client, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	record := &kgo.Record{
		Topic: "orders",
		Value: []byte("hello kafka"),
	}

	ctx := context.Background()
	res, err := client.ProduceSync(ctx, record).First()
	if err != nil {
		log.Fatalf("failed to produce: %v", err)
	}

	fmt.Printf("sent: partition=%d offset=%d\n", res.Partition, res.Offset)
}
