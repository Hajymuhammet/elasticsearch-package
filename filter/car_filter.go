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

type CarFilter struct {
	BrandID           []int64
	ModelID           []int64
	StockID           []int64
	YearMin           *int64
	YearMax           *int64
	PriceMin          *int64
	PriceMax          *int64
	CityID            []int64
	EngineType        []string
	Transmission      []string
	DriveType         []string
	BodyID            []int64
	MileageMin        *int64
	MileageMax        *int64
	EngineCapacityMin *float64
	EngineCapacityMax *float64
	Color             []string
	IsExchange        *bool
	IsCredit          *bool
	PriceOrder        *string // "asc" veya "desc"
	YearOrder         *string // "asc" veya "desc"
	Status            []string
	CreatedAtMin      time.Time
	CreatedAtMax      time.Time
}

func SearchCars(client *elasticsearch.Client, index string, filter *CarFilter) ([]models.Car, error) {
	query := buildESQuery(filter)

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
				Source models.Car `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing response body: %s", err)
	}

	carsList := make([]models.Car, len(r.Hits.Hits))
	for i, hit := range r.Hits.Hits {
		carsList[i] = hit.Source
	}
	return carsList, nil
}

func buildESQuery(filter *CarFilter) map[string]interface{} {
	must := []map[string]interface{}{}

	if len(filter.BrandID) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"brand_id": filter.BrandID}})
	}
	if len(filter.ModelID) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"model_id": filter.ModelID}})
	}
	if len(filter.StockID) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"stock_id": filter.StockID}})
	}
	if filter.YearMin != nil || filter.YearMax != nil {
		r := map[string]interface{}{}
		if filter.YearMin != nil {
			r["gte"] = *filter.YearMin
		}
		if filter.YearMax != nil {
			r["lte"] = *filter.YearMax
		}
		must = append(must, map[string]interface{}{"range": map[string]interface{}{"year": r}})
	}
	if filter.PriceMin != nil || filter.PriceMax != nil {
		r := map[string]interface{}{}
		if filter.PriceMin != nil {
			r["gte"] = *filter.PriceMin
		}
		if filter.PriceMax != nil {
			r["lte"] = *filter.PriceMax
		}
		must = append(must, map[string]interface{}{"range": map[string]interface{}{"price": r}})
	}
	if len(filter.CityID) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"city_id": filter.CityID}})
	}
	if len(filter.EngineType) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"engine_type.keyword": filter.EngineType}})
	}
	if len(filter.Transmission) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"transmission.keyword": filter.Transmission}})
	}
	if len(filter.DriveType) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"drive_type.keyword": filter.DriveType}})
	}
	if len(filter.BodyID) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"body_id": filter.BodyID}})
	}
	if filter.MileageMin != nil || filter.MileageMax != nil {
		r := map[string]interface{}{}
		if filter.MileageMin != nil {
			r["gte"] = *filter.MileageMin
		}
		if filter.MileageMax != nil {
			r["lte"] = *filter.MileageMax
		}
		must = append(must, map[string]interface{}{"range": map[string]interface{}{"mileage": r}})
	}
	if filter.EngineCapacityMin != nil || filter.EngineCapacityMax != nil {
		r := map[string]interface{}{}
		if filter.EngineCapacityMin != nil {
			r["gte"] = *filter.EngineCapacityMin
		}
		if filter.EngineCapacityMax != nil {
			r["lte"] = *filter.EngineCapacityMax
		}
		must = append(must, map[string]interface{}{"range": map[string]interface{}{"engine_capacity": r}})
	}
	if len(filter.Color) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"color.keyword": filter.Color}})
	}
	if filter.IsExchange != nil {
		must = append(must, map[string]interface{}{"term": map[string]interface{}{"is_exchange": *filter.IsExchange}})
	}
	if filter.IsCredit != nil {
		must = append(must, map[string]interface{}{"term": map[string]interface{}{"is_credit": *filter.IsCredit}})
	}
	if len(filter.Status) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"status.keyword": filter.Status}})
	}
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

	// Sort hissəsi
	sort := []map[string]interface{}{}
	if filter.PriceOrder != nil {
		sort = append(sort, map[string]interface{}{"price": map[string]interface{}{"order": *filter.PriceOrder}})
	}
	if filter.YearOrder != nil {
		sort = append(sort, map[string]interface{}{"year": map[string]interface{}{"order": *filter.YearOrder}})
	}

	result := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
	}
	if len(sort) > 0 {
		result["sort"] = sort
	}

	return result
}
