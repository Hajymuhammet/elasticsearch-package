package models

import "time"

type Stock struct {
	ID           int64       `json:"id"`
	UserID       int64       `json:"user_id"`
	UserName     *string     `json:"user_name"`
	PhoneNumber  string      `json:"phone_number"`
	Email        *string     `json:"email"`
	StoreName    *string     `json:"store_name"`
	Images       interface{} `json:"images"` // JSONB
	Logo         interface{} `json:"logo"`   // JSONB
	Address      *string     `json:"address"`
	RegionID     int64       `json:"region_id"`
	CityID       int64       `json:"city_id"`
	CityNameTM   *string     `json:"city_name_tm"`
	CityNameEN   *string     `json:"city_name_en"`
	CityNameRU   *string     `json:"city_name_ru"`
	RegionNameTM *string     `json:"region_name_tm"`
	RegionNameEN *string     `json:"region_name_en"`
	RegionNameRU *string     `json:"region_name_ru"`
	Status       string      `json:"status"`
	Location     Location    `json:"location"`
	Description  *string     `json:"description"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
