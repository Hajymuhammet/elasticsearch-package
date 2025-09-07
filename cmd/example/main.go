package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	esadapter "github.com/salamsites/elasticsearch-package/elasticsearch"
	usecase "github.com/salamsites/elasticsearch-package/service"
)

type Stock struct {
	ID     string  `json:"id" es:"type=keyword"`
	Name   string  `json:"name" es:"type=text"`
	Price  float64 `json:"price" es:"type=double"`
	Amount int     `json:"amount" es:"type:long"`
}

func main() {
	ctx := context.Background()

	client, err := esadapter.NewClient(esadapter.ClientConfig{
		Addresses: []string{"https://localhost:9200"}, // HTTPS ulanylýar
		Username:  "elastic",
		Password:  "+u--nSHZ+G3xfaYxGqzG",
		Timeout:   10 * time.Second,
		TLSConfig: &tls.Config{InsecureSkipVerify: true}, // self-signed sertifikat üçin
	})
	if err != nil {
		log.Fatalf("new client: %v", err)
	}

	repo := esadapter.NewRepository[Stock](client)
	svc := usecase.NewService[Stock](repo)

	// Ensure index
	mapping := svc.BuildMapping(Stock{})
	if err := svc.EnsureIndex(ctx, "stocks", mapping); err != nil {
		log.Fatalf("EnsureIndex failed: %v", err)
	}

	// Index document
	stock := Stock{ID: "1", Name: "Apple", Price: 100.5, Amount: 10}
	if err := svc.Index(ctx, "stocks", stock.ID, stock); err != nil {
		log.Fatalf("Index failed: %v", err)
	}
	fmt.Println("Indexed 1 document")

	// Bulk
	bulkDocs := []Stock{
		{ID: "2", Name: "Google", Price: 2800.0, Amount: 2},
		{ID: "3", Name: "Amazon", Price: 3500.0, Amount: 1},
	}
	if err := svc.BulkIndex(ctx, "stocks", bulkDocs, func(s Stock) string { return s.ID }); err != nil {
		log.Fatalf("BulkIndex failed: %v", err)
	}
	fmt.Println("Bulk indexed")

	// Search
	results, total, err := svc.Search(ctx, "stocks", map[string]any{"match_all": map[string]any{}}, 0, 10, nil)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	fmt.Printf("total=%d results=%v\n", total, results)
}
