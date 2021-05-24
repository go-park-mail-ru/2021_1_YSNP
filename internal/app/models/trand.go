package models

import "time"

//easyjson:json
type UserInterested struct {
	UserID uint64 `json:"userID"`
	Text   string `json:"text"`
}

type UserIArray struct {
	UserID uint64
	Text   []string
}

type Trends struct {
	UserID  uint64
	Popular []Popular
}

type Popular struct {
	Count uint64
	Title string
	Date  time.Time
}

type TrendProducts struct {
	UserID  uint64
	Popular []PopularProduct
}

type PopularProduct struct {
	ProductID uint64
	Time      time.Time
}

type PopularSorter []Popular
type ProductSorter []PopularProduct

func (a PopularSorter) Len() int           { return len(a) }
func (a PopularSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PopularSorter) Less(i, j int) bool { return a[i].Count > a[j].Count }

func (a ProductSorter) Len() int           { return len(a) }
func (a ProductSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProductSorter) Less(i, j int) bool { return a[i].Time.Before(a[j].Time) }
