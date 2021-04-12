package models

type Search struct {
	Category   string `json:"category"`
	Date       string `json:"date"`
	FromAmount int    `json:"fromAmount"`
	Radius     int    `json:"radius"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Search     string `json:"search"`
	Sorting    string `json:"sorting"`
	ToAmount   int    `json:"toAmount"`
}
