package models

import "time"

type Car struct {
	ID             int64       `json:"id"`
	UserId         int64       `json:"user_id"`
	UserName       *string     `json:"user_name"`
	StockId        *int64      `json:"stock_id"`
	StoreName      *string     `json:"store_name"`
	BrandId        int64       `json:"brand_id"`
	BrandName      *string     `json:"brand_name"`
	ModelId        int64       `json:"model_id"`
	ModelName      *string     `json:"model_name"`
	Year           int64       `json:"year"`
	Price          int64       `json:"price"`
	Color          string      `json:"color"`
	Vin            *string     `json:"vin"`
	Description    *string     `json:"description"`
	CityId         int64       `json:"city_id"`
	CityNameTM     *string     `json:"city_name_tm"`
	CityNameEN     *string     `json:"city_name_en"`
	CityNameRU     *string     `json:"city_name_ru"`
	Name           *string     `json:"name"`
	Mail           *string     `json:"mail"`
	PhoneNumber    string      `json:"phone_number"`
	IsComment      bool        `json:"is_comment"`
	IsExchange     bool        `json:"is_exchange"`
	IsCredit       bool        `json:"is_credit"`
	Images         interface{} `json:"images"`
	Status         string      `json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	Mileage        int64       `json:"mileage"`
	EngineCapacity float64     `json:"engine_capacity"`
	EngineType     string      `json:"engine_type"`
	BodyId         int64       `json:"body_id"`
	BodyNameTM     *string     `json:"body_name_tm"`
	BodyNameEN     *string     `json:"body_name_en"`
	BodyNameRU     *string     `json:"body_name_ru"`
	Transmission   string      `json:"transmission"`
	DriveType      string      `json:"drive_type"`
}
