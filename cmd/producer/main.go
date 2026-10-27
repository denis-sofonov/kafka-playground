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

	keys := []string{"user-1", "user-4", "user-8"}

	var records []*kgo.Record
	for i := 0; i < 9; i++ {
		key := keys[i%len(keys)]
		records = append(records, &kgo.Record{
			Topic: "orders",
			Key:   []byte(key),
			Value: []byte(fmt.Sprintf("order-%d", i)),
		})
	}

	ctx := context.Background()
	results := client.ProduceSync(ctx, records...)

	for _, res := range results {
		if res.Err != nil {
			log.Printf("failed to produce: %v", res.Err)
			continue
		}

		fmt.Printf("sent: key=%s value=%s partition=%d offset=%d\n",
			res.Record.Key,
			res.Record.Value,
			res.Record.Partition,
			res.Record.Offset)
	}
}
