package index

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/Hajymuhammet/elasticsearch-package/models"
	"github.com/elastic/go-elasticsearch/v8"
)

var stockMapping = []byte(`
{
  "mappings": {
    "properties": {
      "id": { "type": "long" },
      "user_id": { "type": "long" },
      "user_name": { "type": "keyword" },
      "phone_number": { "type": "keyword" },
      "email": { "type": "keyword" },
      "store_name": { "type": "keyword" },
      "images": { "type": "object", "enabled": true },
      "logo": { "type": "object", "enabled": true },
      "region_id": { "type": "long" },
      "city_id": { "type": "long" },
      "address": { "type": "text" },
      "city_name_tm": { "type": "keyword" },
      "city_name_en": { "type": "keyword" },
      "city_name_ru": { "type": "keyword" },
      "region_name_tm": { "type": "keyword" },
      "region_name_en": { "type": "keyword" },
      "region_name_ru": { "type": "keyword" },
      "status": { "type": "keyword" },
      "description": { "type": "text" },
      "location": {
        "properties": {
          "latitude": { "type": "keyword" },
          "longitude": { "type": "keyword" }
        }
      },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}
`)

func EnsureStockIndex(client *elasticsearch.Client, index string) error {
	res, err := client.Indices.Exists([]string{index})
	if err != nil {
		return fmt.Errorf("error while checking if index exists: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	res, err = client.Indices.Create(
		index,
		client.Indices.Create.WithBody(bytes.NewReader(stockMapping)),
	)
	if err != nil {
		return fmt.Errorf("error creating index %s: %w", index, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create index %s: %s", index, res.String())
	}
	return nil
}

func IndexStock(client *elasticsearch.Client, index string, stock *models.Stock) error {
	data, err := json.Marshal(stock)
	if err != nil {
		return fmt.Errorf("error marshalling stock: %s", err)
	}

	res, err := client.Index(
		index,
		bytes.NewReader(data),
		client.Index.WithDocumentID(fmt.Sprintf("%d", stock.ID)),
		client.Index.WithRefresh("wait_for"),
	)
	if err != nil {
		return fmt.Errorf("error indexing stock ID=%d: %s", stock.ID, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing stock ID=%d: %s", stock.ID, res.String())
	}

	return nil
}

func UpdateStock(client *elasticsearch.Client, index string, stock *models.Stock) error {
	data, err := json.Marshal(map[string]interface{}{
		"doc": stock,
	})

	if err != nil {
		fmt.Println("Error marshalling update stock:", err)
	}

	res, err := client.Update(
		index,
		fmt.Sprintf("%d", stock.ID),
		bytes.NewReader(data),
		client.Update.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error updating document ID=%d: %s", stock.ID, res.String())
	}
	fmt.Printf("Document ID=%d updated successfully\n", stock.ID)
	return nil
}

func DeleteStock(client *elasticsearch.Client, index string, stockID int64) error {
	res, err := client.Delete(
		index,
		fmt.Sprintf("%d", stockID),
		client.Delete.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting document ID=%d: %s", stockID, res.String())
	}

	fmt.Printf("Document ID=%d deleted successfully\n", stockID)

	return nil
}

func BulkUpdateStocks(client *elasticsearch.Client, index string, stocks []models.Stock) error {
	var buf bytes.Buffer

	for _, stock := range stocks {
		meta := []byte(fmt.Sprintf(`{ "update": { "_index": "%s", "_id": "%d" } }%s`, index, stock.ID, "\n"))
		doc, err := json.Marshal(map[string]interface{}{"doc": stock})
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

func BulkDeleteStocks(client *elasticsearch.Client, index string, stockIDs []int64) error {
	var buf bytes.Buffer

	for _, id := range stockIDs {
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
