package filter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Hajymuhammet/elasticsearch-package/models"
	"github.com/elastic/go-elasticsearch/v8"
)

type StockFilter struct {
	UserID       []int64
	CityID       []int64
	RegionID     []int64
	Status       []string
	StoreName    *string
	CityName     *string
	RegionName   *string
	CreatedAtMin time.Time
	CreatedAtMax time.Time
}

func SearchStocks(client *elasticsearch.Client, index string, filter *StockFilter) ([]models.Stock, error) {
	query := buildStockESQuery(filter)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %s", err)
	}

	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex(index),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response: %s", res.String())
	}

	var r struct {
		Hits struct {
			Hits []struct {
				Source models.Stock `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing response body: %s", err)
	}

	stocksList := make([]models.Stock, len(r.Hits.Hits))
	for i, hit := range r.Hits.Hits {
		stocksList[i] = hit.Source
	}

	return stocksList, nil
}

func buildStockESQuery(filter *StockFilter) map[string]interface{} {
	must := []map[string]interface{}{}

	if len(filter.UserID) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"user_id": filter.UserID}})
	}
	if len(filter.CityID) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"city_id": filter.CityID}})
	}
	if len(filter.RegionID) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"region_id": filter.RegionID}})
	}
	if len(filter.Status) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"status.keyword": filter.Status}})
	}
	if filter.StoreName != nil {
		must = append(must, map[string]interface{}{"match": map[string]interface{}{"store_name": *filter.StoreName}})
	}
	if filter.CityName != nil {
		must = append(must, map[string]interface{}{"match": map[string]interface{}{"city_name_tm": *filter.CityName}})
	}
	if filter.RegionName != nil {
		must = append(must, map[string]interface{}{"match": map[string]interface{}{"region_name_tm": *filter.RegionName}})
	}

	// created_at
	if !filter.CreatedAtMin.IsZero() || !filter.CreatedAtMax.IsZero() {
		r := map[string]interface{}{}
		if !filter.CreatedAtMin.IsZero() {
			r["gte"] = filter.CreatedAtMin.Format(time.RFC3339)
		}
		if !filter.CreatedAtMax.IsZero() {
			r["lte"] = filter.CreatedAtMax.Format(time.RFC3339)
		}
		must = append(must, map[string]interface{}{"range": map[string]interface{}{"created_at": r}})
	}

	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
	}
}
