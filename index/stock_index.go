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
      "address": { "type": "text" },
      "region_id": { "type": "long" },
      "city_id": { "type": "long" },
      "city_name_tm": { "type": "keyword" },
      "city_name_en": { "type": "keyword" },
      "city_name_ru": { "type": "keyword" },
      "region_name_tm": { "type": "keyword" },
      "region_name_en": { "type": "keyword" },
      "region_name_ru": { "type": "keyword" },
      "status": { "type": "keyword" },
      "description": { "type": "text" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}
`)

func CreateStockIndex(client *elasticsearch.Client, index string) error {
	res, err := client.Indices.Create(
		index,
		client.Indices.Create.WithBody(bytes.NewReader(stockMapping)),
	)
	if err != nil {
		return err
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
		client.Index.WithRefresh("true"),
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
