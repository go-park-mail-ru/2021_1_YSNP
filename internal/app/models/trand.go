package models


type UserInterested struct {
	UserID uint64 `json:"userID"`
	Text string `json:"text"`
}

type UserIArray struct {
	UserID uint64
	Text []string
}

type Trands struct {
	UserID uint64
	Popular []Popular
}

type Popular struct {
	Count uint64
	Title string
}

type PopularSorter []Popular

func (a PopularSorter) Len() int           { return len(a) }
func (a PopularSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PopularSorter) Less(i, j int) bool { return a[i].Count > a[j].Count }