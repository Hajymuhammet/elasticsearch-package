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
	BrandID            []int64
	ModelID            []int64
	BodyID             []int64
	StockID            []int64
	LoadCapacity       []string
	EngineType         []string
	Transmission       []string
	DriveType          []string
	CityID             []int64
	Color              []string
	BodyType           []string
	CabType            []string
	WheelFormula       []string
	Brakes             []string
	VehicleType        []string
	ForkliftType       []string
	CabSuspension      []string
	SuspensionType     []string
	Status             []string
	Chassis            []string
	BusType            []string
	ExcavatorType      []string
	BulldozerType      []string
	YearMin            *int64
	YearMax            *int64
	PriceMin           *int64
	PriceMax           *int64
	MileageMin         *int64
	MileageMax         *int64
	SeatsMin           *int64
	SeatsMax           *int64
	AxlesMin           *int64
	AxlesMax           *int64
	EngineHoursMin     *int64
	EngineHoursMax     *int64
	LiftingCapacityMin *int64
	LiftingCapacityMax *int64
	EngineCapacityMin  *float64
	EngineCapacityMax  *float64
	Vin                *string
	IsExchange         *bool
	IsCredit           *bool
	PriceOrder         *string // "asc" veya "desc"
	YearOrder          *string // "asc" veya "desc"
	CreatedAtMin       time.Time
	CreatedAtMax       time.Time
	IsCompany          *bool
	IsPrivate          *bool
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

	// Terms filters
	termsFilters := map[string][]string{
		"load_capacity":           filter.LoadCapacity,
		"engine_type":             filter.EngineType,
		"transmission":            filter.Transmission,
		"drive_type":              filter.DriveType,
		"color":                   filter.Color,
		"body_type.keyword":       filter.BodyType,
		"cab_type.keyword":        filter.CabType,
		"wheel_formula.keyword":   filter.WheelFormula,
		"brakes.keyword":          filter.Brakes,
		"vehicle_type.keyword":    filter.VehicleType,
		"forklift_type.keyword":   filter.ForkliftType,
		"cab_suspension.keyword":  filter.CabSuspension,
		"suspension_type.keyword": filter.SuspensionType,
		"status.keyword":          filter.Status,
		"chassis.keyword":         filter.Chassis,
		"bus_type.keyword":        filter.BusType,
		"excavator_type.keyword":  filter.ExcavatorType,
		"bulldozer_type.keyword":  filter.BulldozerType,
	}

	for field, values := range termsFilters {
		if len(values) > 0 {
			must = append(must, map[string]interface{}{"terms": map[string]interface{}{field: values}})
		}
	}

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
	if len(filter.CityID) > 0 {
		must = append(must, map[string]interface{}{"terms": map[string]interface{}{"city_id": filter.CityID}})
	}
	if filter.Vin != nil {
		must = append(must, map[string]interface{}{"term": map[string]interface{}{"vin.keyword": *filter.Vin}})
	}
	if filter.IsExchange != nil {
		must = append(must, map[string]interface{}{"term": map[string]interface{}{"is_exchange": *filter.IsExchange}})
	}
	if filter.IsCredit != nil {
		must = append(must, map[string]interface{}{"term": map[string]interface{}{"is_credit": *filter.IsCredit}})
	}

	// Range filters
	rangeFields := map[string][2]*interface{}{
		"year":             {toIface(filter.YearMin), toIface(filter.YearMax)},
		"price":            {toIface(filter.PriceMin), toIface(filter.PriceMax)},
		"mileage":          {toIface(filter.MileageMin), toIface(filter.MileageMax)},
		"seats":            {toIface(filter.SeatsMin), toIface(filter.SeatsMax)},
		"axles":            {toIface(filter.AxlesMin), toIface(filter.AxlesMax)},
		"engine_hours":     {toIface(filter.EngineHoursMin), toIface(filter.EngineHoursMax)},
		"lifting_capacity": {toIface(filter.LiftingCapacityMin), toIface(filter.LiftingCapacityMax)},
		"engine_capacity":  {toIface(filter.EngineCapacityMin), toIface(filter.EngineCapacityMax)},
	}

	for field, bounds := range rangeFields {
		r := map[string]interface{}{}
		if bounds[0] != nil {
			r["gte"] = *bounds[0]
		}
		if bounds[1] != nil {
			r["lte"] = *bounds[1]
		}
		if len(r) > 0 {
			must = append(must, map[string]interface{}{"range": map[string]interface{}{field: r}})
		}
	}

	// CreatedAt range
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
		if *filter.IsCompany && *filter.IsPrivate {
		} else if *filter.IsCompany {
			must = append(must, map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						{"range": map[string]interface{}{"stock_id": map[string]interface{}{"gt": 0}}},
						{"exists": map[string]interface{}{"field": "stock_id"}},
					},
				},
			})
		} else if *filter.IsPrivate {
			must = append(must, map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						{"term": map[string]interface{}{"stock_id": 0}},
						{"bool": map[string]interface{}{"must_not": map[string]interface{}{"exists": map[string]interface{}{"field": "stock_id"}}}},
					},
				},
			})
		}
	} else if filter.IsCompany != nil && *filter.IsCompany {
		must = append(must, map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{"range": map[string]interface{}{"stock_id": map[string]interface{}{"gt": 0}}},
					{"exists": map[string]interface{}{"field": "stock_id"}},
				},
			},
		})
	} else if filter.IsPrivate != nil && *filter.IsPrivate {
		must = append(must, map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{"term": map[string]interface{}{"stock_id": 0}},
					{"bool": map[string]interface{}{"must_not": map[string]interface{}{"exists": map[string]interface{}{"field": "stock_id"}}}},
				},
			},
		})
	}

	// Sort
	sort := []map[string]interface{}{}
	if filter.PriceOrder != nil {
		sort = append(sort, map[string]interface{}{"price": map[string]interface{}{"order": *filter.PriceOrder}})
	}
	if filter.YearOrder != nil {
		sort = append(sort, map[string]interface{}{"year": map[string]interface{}{"order": *filter.YearOrder}})
	}

	result := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{"must": must},
		},
	}
	if len(sort) > 0 {
		result["sort"] = sort
	}

	return result
}

// helper: safely convert typed pointers to *interface{}
func toIface[T any](v *T) *interface{} {
	if v == nil {
		return nil
	}
	var i interface{} = *v
	return &i
}
