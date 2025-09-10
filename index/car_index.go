package index

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/Hajymuhammet/elasticsearch-package/models"
	"github.com/elastic/go-elasticsearch/v8"
)

var carMapping = []byte(`
{
  "mappings": {
    "properties": {
      "id": { "type": "long" },
      "user_id": { "type": "long" },
      "user_name": { "type": "keyword" },
      "stock_id": { "type": "long" },
      "store_name": { "type": "keyword" },
      "brand_id": { "type": "long" },
      "brand_name": { "type": "keyword" },
      "model_id": { "type": "long" },
      "model_name": { "type": "keyword" },
      "year": { "type": "long" },
      "price": { "type": "long" },
      "color": { "type": "keyword" },
      "vin": { "type": "keyword" },
      "description": { "type": "text" },
      "city_id": { "type": "long" },
      "city_name_tm": { "type": "keyword" },
      "city_name_en": { "type": "keyword" },
      "city_name_ru": { "type": "keyword" },
      "name": { "type": "keyword" },
      "mail": { "type": "keyword" },
      "phone_number": { "type": "keyword" },
      "is_comment": { "type": "boolean" },
      "is_exchange": { "type": "boolean" },
      "is_credit": { "type": "boolean" },
      "images": { "type": "object", "enabled": true },
      "status": { "type": "keyword" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" },
      "mileage": { "type": "long" },
      "engine_capacity": { "type": "double" },
      "engine_type": { "type": "keyword" },
      "body_id": { "type": "long" },
      "body_name_tm": { "type": "keyword" },
      "body_name_en": { "type": "keyword" },
      "body_name_ru": { "type": "keyword" },
      "transmission": { "type": "keyword" },
      "drive_type": { "type": "keyword" },
      "options": { "type": "long" }
    }
  }
}
`)

// CreateCarIndex göni index-i mapping bilen döretmek üçin
func EnsureCarIndex(client *elasticsearch.Client, index string) error {
	res, err := client.Indices.Exists([]string{index})
	if err != nil {
		return fmt.Errorf("Error checking if index exists: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		fmt.Println("Using existing index:", index)
		return nil
	}

	res, err = client.Indices.Create(
		index,
		client.Indices.Create.WithBody(bytes.NewReader(carMapping)),
	)
	if err != nil {
		return fmt.Errorf("Error creating index %s: %v", index, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create index %s: %s", index, res.String())
	}

	fmt.Println("Index created successfully:", index)

	return nil
}

// IndexCar dokumenti goşmak
func IndexCar(client *elasticsearch.Client, index string, car *models.Car) error {
	data, err := json.Marshal(car)
	if err != nil {
		return err
	}

	res, err := client.Index(
		index,
		bytes.NewReader(data),
		client.Index.WithDocumentID(fmt.Sprintf("%d", car.ID)),
		client.Index.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document ID=%d: %s", car.ID, res.String())
	}
	fmt.Printf("Document ID=%d indexed successfully\n", car.ID)

	return nil
}
