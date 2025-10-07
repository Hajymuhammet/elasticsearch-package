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

type MotoFilter struct {
	BrandID             []int64
	ModelID             []int64
	BodyID              []int64
	StockID             []int64
	YearMin             *int32
	YearMax             *int32
	PriceMin            *int64
	PriceMax            *int64
	CityID              []int64
	EngineType          []string
	TypeMotorcycles     []string
	MileageMin          *int64
	MileageMax          *int64
	VolumeMin           *int64
	VolumeMax           *int64
	Color               []string
	IsExchange          *bool
	IsCredit            *bool
	Status              []string
	NumberOfClockCycles []int64
	AirType             []string
	Options             []int64
	PriceOrder          *string // "asc" / "desc"
	YearOrder           *string // "asc" / "desc"
	CreatedAtMin        time.Time
	CreatedAtMax        time.Time
	IsCompany           *bool
	IsPrivate           *bool
}

func SearchMotos(client *elasticsearch.Client, index string, filter *MotoFilter) ([]models.Moto, error) {
	query := buildMotoESQuery(filter)

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
				Source models.Moto `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing response body: %s", err)
	}

	motosList := make([]models.Moto, len(r.Hits.Hits))
	for i, hit := range r.Hits.Hits {
		motosList[i] = hit.Source
	}

	return motosList, nil
}

func buildMotoESQuery(filter *MotoFilter) map[string]interface{} {
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
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"engine_type": filter.EngineType}})
	}
	if len(filter.TypeMotorcycles) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"type_motorcycles": filter.TypeMotorcycles}})
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
	if filter.VolumeMin != nil || filter.VolumeMax != nil {
		r := map[string]interface{}{}
		if filter.VolumeMin != nil {
			r["gte"] = *filter.VolumeMin
		}
		if filter.VolumeMax != nil {
			r["lte"] = *filter.VolumeMax
		}
		must = append(must, map[string]interface{}{"range": map[string]interface{}{"volume": r}})
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
	if len(filter.NumberOfClockCycles) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"number_of_clock_cycles": filter.NumberOfClockCycles}})
	}
	if len(filter.AirType) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"air_type.keyword": filter.AirType}})
	}
	if len(filter.Options) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"options": filter.Options}})
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

	if filter.IsCompany != nil && filter.IsPrivate != nil {
		if !*filter.IsCompany && !*filter.IsPrivate {
		} else if *filter.IsCompany {
			must = append(must, map[string]interface{}{
				"bool": map[string]interface{}{
					"must_not": map[string]interface{}{
						"term": map[string]interface{}{"stock_id": 0},
					},
				},
			})
		} else if *filter.IsPrivate {
			must = append(must, map[string]interface{}{
				"term": map[string]interface{}{"stock_id": 0},
			})
		}
	} else if filter.IsCompany != nil && *filter.IsCompany {
		must = append(must, map[string]interface{}{
			"bool": map[string]interface{}{
				"must_not": map[string]interface{}{
					"term": map[string]interface{}{"stock_id": 0},
				},
			},
		})
	} else if filter.IsPrivate != nil && *filter.IsPrivate {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{"stock_id": 0},
		})
	}

	// Sorting
	sort := []map[string]interface{}{}
	if filter.PriceOrder != nil {
		sort = append(sort, map[string]interface{}{"price": map[string]interface{}{"order": *filter.PriceOrder}})
	}
	if filter.YearOrder != nil {
		sort = append(sort, map[string]interface{}{"year": map[string]interface{}{"order": *filter.YearOrder}})
	}
	if len(sort) == 0 {
		// default sort by created_at desc
		sort = append(sort, map[string]interface{}{"created_at": map[string]interface{}{"order": "desc"}})
	}

	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{"must": must},
		},
		"sort": sort,
	}
}
