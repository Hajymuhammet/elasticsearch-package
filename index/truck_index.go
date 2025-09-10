package index

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/Hajymuhammet/elasticsearch-package/models"
	"github.com/elastic/go-elasticsearch/v8"
)

var truckMapping = []byte(`
{
  "mappings": {
    "properties": {
      "id": { "type": "long" },
      "user_id": { "type": "long" },
      "user_name": { "type": "keyword" },
      "stock_id": { "type": "long" },
      "store_name": { "type": "keyword" },
      "body_id": { "type": "long" },
      "body_name_tm": { "type": "keyword" },
      "body_name_en": { "type": "keyword" },
      "body_name_ru": { "type": "keyword" },
      "brand_id": { "type": "long" },
      "brand_name": { "type": "keyword" },
      "model_id": { "type": "long" },
      "model_name": { "type": "keyword" },
      "load_capacity": { "type": "double" },
      "price": { "type": "long" },
      "body_type": { "type": "keyword" },
      "drive_type": { "type": "keyword" },
      "transmission": { "type": "keyword" },
      "engine_type": { "type": "keyword" },
      "year": { "type": "long" },
      "seats": { "type": "long" },
      "cab_type": { "type": "keyword" },
      "wheel_formula": { "type": "keyword" },
      "chassis": { "type": "keyword" },
      "cab_suspension": { "type": "keyword" },
      "bus_type": { "type": "keyword" },
      "suspension_type": { "type": "keyword" },
      "brakes": { "type": "keyword" },
      "axles": { "type": "long" },
      "engine_hours": { "type": "long" },
      "vehicle_type": { "type": "keyword" },
      "engine_capacity": { "type": "double" },
      "forklift_type": { "type": "keyword" },
      "lifting_capacity": { "type": "long" },
      "mileage": { "type": "long" },
      "excavator_type": { "type": "keyword" },
      "bulldozer_type": { "type": "keyword" },
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
      "options": { "type": "long" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}
`)

func EnsureTruckIndex(client *elasticsearch.Client, indexName string) error {
	res, err := client.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("error checking if index exists: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	res, err = client.Indices.Create(
		indexName,
		client.Indices.Create.WithBody(bytes.NewReader(truckMapping)),
	)
	if err != nil {
		return fmt.Errorf("error while creating truck index: %s: %w", indexName, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create truck index %s: %s", indexName, res.String())
	}

	fmt.Println("Created truck index:", indexName)

	return nil
}

func IndexTruck(client *elasticsearch.Client, indexName string, truck *models.Truck) error {
	data, err := json.Marshal(truck)
	if err != nil {
		fmt.Println("Error marshalling truck index:", err)
	}

	res, err := client.Index(
		indexName,
		bytes.NewReader(data),
		client.Index.WithDocumentID(fmt.Sprintf("%d", truck.Id)),
		client.Index.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing truck ID=%d: %s", truck.Id, res.String())
	}

	fmt.Printf("Truck ID=%d indexed successfully\n", truck.Id)
	return nil
}

func UpdateTruck(client *elasticsearch.Client, index string, truck *models.Truck) error {
	data, err := json.Marshal(map[string]interface{}{"doc": truck})
	if err != nil {
		return err
	}

	res, err := client.Update(
		index,
		fmt.Sprintf("%d", truck.Id),
		bytes.NewReader(data),
		client.Update.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error updating truck ID=%d: %s", truck.Id, res.String())
	}

	fmt.Printf("Truck ID=%d updated successfully\n", truck.Id)
	return nil
}

func DeleteTruck(client *elasticsearch.Client, index string, truckID int64) error {
	res, err := client.Delete(
		index,
		fmt.Sprintf("%d", truckID),
		client.Delete.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting truck ID=%d: %s", truckID, res.String())
	}

	fmt.Printf("Truck ID=%d deleted successfully\n", truckID)
	return nil
}

func BulkUpdateTrucks(client *elasticsearch.Client, index string, trucks []models.Truck) error {
	var buf bytes.Buffer

	for _, truck := range trucks {
		meta := []byte(fmt.Sprintf(`{ "update": { "_index": "%s", "_id": "%d" } }%s`, index, truck.Id, "\n"))
		doc, err := json.Marshal(map[string]interface{}{"doc": truck})
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

func BulkDeleteTrucks(client *elasticsearch.Client, index string, truckIDs []int64) error {
	var buf bytes.Buffer

	for _, id := range truckIDs {
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
