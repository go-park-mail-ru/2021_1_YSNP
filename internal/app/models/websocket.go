package models

import "encoding/json"

type WSMessageReq struct {
	UserID uint64  `json:"-"`
	Type string `json:"type"`
	Data CustomData `json:"data"`
}

type CustomData struct {
	RequestData json.RawMessage `json:"request_data"`
	TypeData json.RawMessage `json:"type_data"`
}

type WSMessageResp struct {
	UserID uint64 `json:"-"`
	Status int  `json:"status"`
	Type string  `json:"type"`
	Data json.RawMessage `json:"data"`
}
