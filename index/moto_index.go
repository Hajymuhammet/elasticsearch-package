package index

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/Hajymuhammet/elasticsearch-package/models"
	"github.com/elastic/go-elasticsearch/v8"
)

var motoMapping = []byte(`
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
      "type_motorcycles": { "type": "keyword" },
      "year": { "type": "long" },
      "price": { "type": "long" },
      "volume": { "type": "long" },
      "engine_type": { "type": "keyword" },
      "number_of_clock_cycles": { "type": "long" },
      "mileage": { "type": "long" },
      "air_type": { "type": "keyword" },
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
      "updated_at": { "type": "date" }
    }
  }
}
`)

func EnsureMotoIndex(client *elasticsearch.Client, index string) error {
	res, err := client.Indices.Exists([]string{index})
	if err != nil {
		return fmt.Errorf("error checking index: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		fmt.Println("Using existing moto index:", index)
		return nil
	}

	res, err = client.Indices.Create(
		index,
		client.Indices.Create.WithBody(bytes.NewReader(motoMapping)),
	)
	if err != nil {
		return fmt.Errorf("error creating moto index %s: %v", index, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create moto index %s: %s", index, res.String())
	}

	fmt.Println("Moto index created successfully:", index)
	return nil
}

func IndexMoto(client *elasticsearch.Client, index string, moto *models.Moto) error {
	data, err := json.Marshal(moto)
	if err != nil {
		return err
	}

	res, err := client.Index(
		index,
		bytes.NewReader(data),
		client.Index.WithDocumentID(fmt.Sprintf("%d", moto.Id)),
		client.Index.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing moto ID=%d: %s", moto.Id, res.String())
	}

	fmt.Printf("Moto ID=%d indexed successfully\n", moto.Id)
	return nil
}

func UpdateMoto(client *elasticsearch.Client, index string, moto *models.Moto) error {
	data, err := json.Marshal(map[string]interface{}{
		"doc":           moto,
		"doc_as_upsert": true,
	})
	if err != nil {
		return err
	}

	res, err := client.Update(
		index,
		fmt.Sprintf("%d", moto.Id),
		bytes.NewReader(data),
		client.Update.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error updating moto ID=%d: %s", moto.Id, res.String())
	}

	fmt.Printf("Moto ID=%d updated successfully\n", moto.Id)
	return nil
}

func DeleteMoto(client *elasticsearch.Client, index string, motoID int64) error {
	res, err := client.Delete(
		index,
		fmt.Sprintf("%d", motoID),
		client.Delete.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting moto ID=%d: %s", motoID, res.String())
	}

	fmt.Printf("Moto ID=%d deleted successfully\n", motoID)
	return nil
}

func BulkUpdateMotos(client *elasticsearch.Client, index string, motos []models.Moto) error {
	var buf bytes.Buffer

	for _, moto := range motos {
		meta := []byte(fmt.Sprintf(`{ "update": { "_index": "%s", "_id": "%d" } }%s`, index, moto.Id, "\n"))
		doc, err := json.Marshal(map[string]interface{}{"doc": moto})
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

func BulkDeleteMotos(client *elasticsearch.Client, index string, motoIDs []int64) error {
	var buf bytes.Buffer

	for _, id := range motoIDs {
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
