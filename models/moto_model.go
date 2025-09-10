package models

import "time"

type Moto struct {
	Id                  int64       `json:"id"`
	UserId              *int64      `json:"user_id"`
	UserName            *string     `json:"user_name"`
	StockId             *int64      `json:"stock_id"`
	StoreName           *string     `json:"store_name"`
	BodyId              int64       `json:"body_id"`
	BodyNameTM          *string     `json:"body_name_tm"`
	BodyNameEN          *string     `json:"body_name_en"`
	BodyNameRU          *string     `json:"body_name_ru"`
	BrandId             int64       `json:"brand_id"`
	BrandName           *string     `json:"brand_name"`
	ModelId             int64       `json:"model_id"`
	ModelName           *string     `json:"model_name"`
	TypeMotorcycles     *string     `json:"type_motorcycles"`
	Year                int32       `json:"year"`
	Price               int64       `json:"price"`
	Volume              int64       `json:"volume"`
	EngineType          *string     `json:"engine_type"`
	NumberOfClockCycles *int32      `json:"number_of_clock_cycles"`
	Mileage             *int64      `json:"mileage"`
	AirType             *string     `json:"air_type"`
	Color               string      `json:"color"`
	Vin                 *string     `json:"vin"`
	Description         *string     `json:"description"`
	CityId              *int64      `json:"city_id"`
	CityNameTM          *string     `json:"city_name_tm"`
	CityNameEN          *string     `json:"city_name_en"`
	CityNameRU          *string     `json:"city_name_ru"`
	Name                *string     `json:"name"`
	Mail                *string     `json:"mail"`
	PhoneNumber         string      `json:"phone_number"`
	Options             []int64     `json:"options"`
	IsComment           bool        `json:"is_comment"`
	IsExchange          bool        `json:"is_exchange"`
	IsCredit            bool        `json:"is_credit"`
	Images              interface{} `json:"images"`
	Status              string      `json:"status"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
}
