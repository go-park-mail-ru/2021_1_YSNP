package models

type Search struct {
	Category   string  `valid:"type(string)"`
	Date       string  `valid:"type(string)"`
	FromAmount int     `valid:"numeric"`
	ToAmount   int     `valid:"numeric"`
	Radius     uint64  `valid:"numeric"`
	Latitude   float64 `valid:"latitude"`
	Longitude  float64 `valid:"longitude"`
	Search     string  `valid:"type(string)"`
	Sorting    string  `valid:"type(string)"`
	From       uint64  `valid:"numeric"`
	Count      uint64  `valid:"numeric"`
}
