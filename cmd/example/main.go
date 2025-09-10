package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	esadapter "github.com/Hajymuhammet/elasticsearch-package/clients"
	"github.com/Hajymuhammet/elasticsearch-package/filter"
	"github.com/Hajymuhammet/elasticsearch-package/index"
	"github.com/Hajymuhammet/elasticsearch-package/models"
)

const (
	carIndexName = "cars"
)

func main() {
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

	if err := index.EnsureCarIndex(es, carIndexName); err != nil {
		log.Fatalf("failed to ensure index: %v", err)
	}

	car := &models.Car{
		ID:             time.Now().UnixNano(),
		UserId:         101,
		UserName:       ptrString("John Doe"),
		StockId:        ptrInt64(201),
		StoreName:      ptrString("AutoStore"),
		BrandId:        10,
		BrandName:      ptrString("Toyota"),
		ModelId:        1001,
		ModelName:      ptrString("Highlender"),
		Year:           2018,
		Price:          15000,
		Color:          "blue",
		Vin:            ptrString("1HGBH41JXMN109186"),
		Description:    ptrString("Well maintained car"),
		CityId:         1,
		CityNameTM:     ptrString("Ashgabat"),
		CityNameEN:     ptrString("Ashgabat"),
		CityNameRU:     ptrString("Ашхабад"),
		Name:           ptrString("John Car"),
		Mail:           ptrString("john@example.com"),
		PhoneNumber:    "+99361234567",
		IsComment:      true,
		IsExchange:     false,
		IsCredit:       true,
		Images:         []string{"img1.jpg", "img2.jpg"},
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Mileage:        50000,
		EngineCapacity: 1.8,
		EngineType:     "Of",
		BodyId:         1,
		BodyNameTM:     ptrString("Sedan"),
		BodyNameEN:     ptrString("Sedan"),
		BodyNameRU:     ptrString("Седан"),
		Transmission:   "Mehanic",
		DriveType:      "AWD",
		Options:        []int64{1, 2, 3},
	}

	if err := index.IndexCar(es, carIndexName, car); err != nil {
		fmt.Println("index error:", err)
	} else {
		fmt.Println("Document indexed successfully")
	}

	// Filter bilen sorag
	carFilter := &filter.CarFilter{
		Color: []string{"blue"},
	}

	cars, err := filter.SearchCars(es, carIndexName, carFilter)
	if err != nil {
		fmt.Println("search error:", err)
	} else {
		fmt.Printf("Found %d cars:\n", len(cars))
		data, _ := json.MarshalIndent(cars, "", "  ")
		fmt.Println(string(data))
	}
}

func ptrString(s string) *string { return &s }
func ptrInt64(i int64) *int64    { return &i }
