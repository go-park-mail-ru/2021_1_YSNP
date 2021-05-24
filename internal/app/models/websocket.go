package models

import "encoding/json"

//easyjson:json
type WSMessageReq struct {
	UserID uint64     `json:"-"`
	Type   string     `json:"type"`
	Data   CustomData `json:"data"`
}

//easyjson:json
type CustomData struct {
	RequestData json.RawMessage `json:"request_data"`
	TypeData    json.RawMessage `json:"type_data"`
}

//easyjson:json
type WSMessageResp struct {
	UserID uint64          `json:"-"`
	Status int             `json:"status"`
	Type   string          `json:"type"`
	Data   json.RawMessage `json:"data"`
}
