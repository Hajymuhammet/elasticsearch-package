package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	esadapter "github.com/Hajymuhammet/elasticsearch-package/clients"
	"github.com/Hajymuhammet/elasticsearch-package/index"
	"github.com/Hajymuhammet/elasticsearch-package/models"
)

const (
	carIndexName = "cars"
)

func ptrString(s string) *string { return &s }
func ptrInt64(i int64) *int64    { return &i }

func main() {
	// Elasticsearch client
	es, err := esadapter.NewClient(esadapter.ClientConfig{
		Addresses: []string{"http://10.192.1.127:9200"},
		Username:  "elastic",
		Password:  "1Z*Ywl6qbw0PXqoeKzJ3",
		Timeout:   10 * time.Second,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	if err != nil {
		log.Fatalf("new client: %v", err)
	}

	// Indeksi barla ýa-da döret
	if err := index.EnsureCarIndex(es, carIndexName); err != nil {
		log.Fatalf("failed to ensure index: %v", err)
	}

	cars := []models.Car{
		{ID: 1, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(9), StoreName: ptrString("test"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2020, Price: 20000, Color: "red", PhoneNumber: "+99360000001", Status: "accepted", CityId: 3, CityNameTM: ptrString("turkmenbasy"), CityNameEN: ptrString("turkmenbasy"), CityNameRU: ptrString("turkmenbasy"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(8), StoreName: ptrString("tetsAuto"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2019, Price: 18000, Color: "blue", PhoneNumber: "+99360000002", Status: "accepted", CityId: 3, CityNameTM: ptrString("turkmenbasy"), CityNameEN: ptrString("turkmenbasy"), CityNameRU: ptrString("turkmenbasy"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 3, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(7), StoreName: ptrString("Awtomobile"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2021, Price: 25000, Color: "green", PhoneNumber: "+99360000003", Status: "accepted", CityId: 3, CityNameTM: ptrString("turkmenbasy"), CityNameEN: ptrString("turkmenbasy"), CityNameRU: ptrString("turkmenbasy"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 4, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(6), StoreName: ptrString("StoreAsgabat"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2018, Price: 15000, Color: "yellow", PhoneNumber: "+99360000004", Status: "accepted", CityId: 3, CityNameTM: ptrString("turkmenbasy"), CityNameEN: ptrString("turkmenbasy"), CityNameRU: ptrString("turkmenbasy"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 5, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(5), StoreName: ptrString("AsgabatStore"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2022, Price: 30000, Color: "black", PhoneNumber: "+99360000005", Status: "accepted", CityId: 2, CityNameTM: ptrString("bayramaly"), CityNameEN: ptrString("bayramaly"), CityNameRU: ptrString("bayramaly"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 6, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(4), StoreName: ptrString("AsgabatStore"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2017, Price: 12000, Color: "white", PhoneNumber: "+99360000006", Status: "accepted", CityId: 2, CityNameTM: ptrString("bayramaly"), CityNameEN: ptrString("bayramaly"), CityNameRU: ptrString("bayramaly"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 7, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(3), StoreName: ptrString("AsgabatStore"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2016, Price: 10000, Color: "gray", PhoneNumber: "+99360000007", Status: "accepted", CityId: 2, CityNameTM: ptrString("bayramaly"), CityNameEN: ptrString("bayramaly"), CityNameRU: ptrString("bayramaly"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 8, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(2), StoreName: ptrString("AsgabatStore"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2020, Price: 22000, Color: "blue", PhoneNumber: "+99360000008", Status: "accepted", CityId: 2, CityNameTM: ptrString("bayramaly"), CityNameEN: ptrString("bayramaly"), CityNameRU: ptrString("bayramaly"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 9, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(9), StoreName: ptrString("test"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2019, Price: 18000, Color: "red", PhoneNumber: "+99360000009", Status: "accepted", CityId: 3, CityNameTM: ptrString("turkmenbasy"), CityNameEN: ptrString("turkmenbasy"), CityNameRU: ptrString("turkmenbasy"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 10, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(8), StoreName: ptrString("tetsAuto"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2018, Price: 16000, Color: "green", PhoneNumber: "+99360000010", Status: "accepted", CityId: 3, CityNameTM: ptrString("turkmenbasy"), CityNameEN: ptrString("turkmenbasy"), CityNameRU: ptrString("turkmenbasy"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 11, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(7), StoreName: ptrString("Awtomobile"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2021, Price: 27000, Color: "yellow", PhoneNumber: "+99360000011", Status: "accepted", CityId: 3, CityNameTM: ptrString("turkmenbasy"), CityNameEN: ptrString("turkmenbasy"), CityNameRU: ptrString("turkmenbasy"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 12, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(6), StoreName: ptrString("StoreAsgabat"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2015, Price: 9000, Color: "white", PhoneNumber: "+99360000012", Status: "accepted", CityId: 3, CityNameTM: ptrString("turkmenbasy"), CityNameEN: ptrString("turkmenbasy"), CityNameRU: ptrString("turkmenbasy"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 13, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(5), StoreName: ptrString("AsgabatStore"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2022, Price: 31000, Color: "black", PhoneNumber: "+99360000013", Status: "accepted", CityId: 2, CityNameTM: ptrString("bayramaly"), CityNameEN: ptrString("bayramaly"), CityNameRU: ptrString("bayramaly"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 14, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(4), StoreName: ptrString("AsgabatStore"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2017, Price: 13000, Color: "gray", PhoneNumber: "+99360000014", Status: "accepted", CityId: 2, CityNameTM: ptrString("bayramaly"), CityNameEN: ptrString("bayramaly"), CityNameRU: ptrString("bayramaly"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 15, UserId: 1, UserName: ptrString(""), StockId: ptrInt64(3), StoreName: ptrString("AsgabatStore"), BrandId: 5, BrandName: ptrString("dededed"), ModelId: 3, ModelName: ptrString("TEST"), Year: 2016, Price: 11000, Color: "blue", PhoneNumber: "+99360000015", Status: "accepted", CityId: 2, CityNameTM: ptrString("bayramaly"), CityNameEN: ptrString("bayramaly"), CityNameRU: ptrString("bayramaly"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	// Indeksle we terminala çykar
	for _, car := range cars {
		if err := index.IndexCar(es, carIndexName, &car); err != nil {
			fmt.Println("index error:", err)
		} else {
			fmt.Printf("Document ID=%d indexed successfully\n", car.ID)
		}
	}
}
