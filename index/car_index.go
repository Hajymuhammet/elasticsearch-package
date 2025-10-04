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
      "drive_type": { "type": "keyword" }
    }
  }
}
`)

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

// IndexCar
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

func UpdateCar(client *elasticsearch.Client, index string, car *models.Car) error {
	data, err := json.Marshal(map[string]interface{}{
		"doc":           car,
		"doc_as_upsert": true,
	})

	if err != nil {
		fmt.Println("Error marshalling update car:", err)
	}

	res, err := client.Update(
		index,
		fmt.Sprintf("%d", car.ID),
		bytes.NewReader(data),
		client.Update.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error updating document ID=%d: %s", car.ID, res.String())
	}
	fmt.Printf("Document ID=%d updated successfully\n", car.ID)
	return nil
}

func DeleteCar(client *elasticsearch.Client, index string, carID int64) error {
	res, err := client.Delete(
		index,
		fmt.Sprintf("%d", carID),
		client.Delete.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting document ID=%d: %s", carID, res.String())
	}

	fmt.Printf("Document ID=%d deleted successfully\n", carID)

	return nil
}

func BulkUpdateCars(client *elasticsearch.Client, index string, cars []models.Car) error {
	var buf bytes.Buffer

	for _, car := range cars {
		meta := []byte(fmt.Sprintf(`{ "update": { "_index": "%s", "_id": "%d" } }%s`, index, car.ID, "\n"))
		doc, err := json.Marshal(map[string]interface{}{"doc": car})
		if err != nil {
			return err
		}
		doc = append(doc, "\n"...)
		buf.Write(meta)
		buf.Write(doc)
	}

	res, err := client.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk update error: %s", res.String())
	}

	fmt.Println("Bulk update successful")
	return nil
}

func BulkDeleteCars(client *elasticsearch.Client, index string, carIDs []int64) error {
	var buf bytes.Buffer

	for _, id := range carIDs {
		meta := []byte(fmt.Sprintf(`{ "delete": { "_index": "%s", "_id": "%d" } }%s`, index, id, "\n"))
		buf.Write(meta)
	}

	res, err := client.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk delete error: %s", res.String())
	}

	fmt.Println("Bulk delete successful")
	return nil
}

func DeleteFeedsByUserID(client *elasticsearch.Client, index string, userID int64) error {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"user_id": userID,
			},
		},
	}

	data, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("Error marshalling query: %w", err)
	}

	res, err := client.DeleteByQuery(
		[]string{index},
		bytes.NewReader(data),
		client.DeleteByQuery.WithRefresh(true),
	)
	if err != nil {
		return fmt.Errorf("error deleting by user_id=%d: %v", userID, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("delete by query failed for user_id=%d: %s", userID, res.String())
	}

	fmt.Printf("All documents with user_id=%d deleted successfully from index=%s\n", userID, index)
	return nil
}
