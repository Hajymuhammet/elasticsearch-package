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

type TruckFilter struct {
	BrandID         []int64
	ModelID         []int64
	BodyID          []int64
	YearMin         *int64
	YearMax         *int64
	PriceMin        *int64
	PriceMax        *int64
	CityID          []int64
	EngineType      []string
	Transmission    []string
	DriveType       []string
	LoadCapacityMin *float64
	LoadCapacityMax *float64
	MileageMin      *int64
	MileageMax      *int64
	EngineCapacity  *float64
	Color           []string
	IsExchange      *bool
	IsCredit        *bool
	Status          []string
	CreatedAtMin    time.Time
	CreatedAtMax    time.Time
}

func SearchTrucks(client *elasticsearch.Client, index string, filter *TruckFilter) ([]models.Truck, error) {
	query := buildTruckESQuery(filter)

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
				Source models.Truck `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing response body: %s", err)
	}

	trucksList := make([]models.Truck, len(r.Hits.Hits))
	for i, hit := range r.Hits.Hits {
		trucksList[i] = hit.Source
	}

	return trucksList, nil
}

func buildTruckESQuery(filter *TruckFilter) map[string]interface{} {
	must := []map[string]interface{}{}

	if len(filter.BrandID) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"brand_id": filter.BrandID}})
	}
	if len(filter.ModelID) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"model_id": filter.ModelID}})
	}
	if len(filter.BodyID) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"body_id": filter.BodyID}})
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
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"engine_type": filter.EngineType}})
	}
	if len(filter.Transmission) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"transmission": filter.Transmission}})
	}
	if len(filter.DriveType) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"drive_type": filter.DriveType}})
	}
	if filter.LoadCapacityMin != nil || filter.LoadCapacityMax != nil {
		r := map[string]interface{}{}
		if filter.LoadCapacityMin != nil {
			r["gte"] = *filter.LoadCapacityMin
		}
		if filter.LoadCapacityMax != nil {
			r["lte"] = *filter.LoadCapacityMax
		}
		must = append(must, map[string]interface{}{"range": map[string]interface{}{"load_capacity": r}})
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
	if filter.EngineCapacity != nil {
		must = append(must, map[string]interface{}{"term": map[string]interface{}{"engine_capacity": *filter.EngineCapacity}})
	}
	if len(filter.Color) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"color": filter.Color}})
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

	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{"must": must},
		},
	}
}
