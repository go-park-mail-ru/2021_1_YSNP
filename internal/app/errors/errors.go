package errors

import "encoding/json"

type Error struct {
	Message string `json:"message"`
}

type Success struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	panic("implement me")
}

func JSONError(message string) []byte {
	jsonError, err := json.Marshal(Error{Message: message})
	if err != nil {
		return []byte("")
	}
	return jsonError
}

func JSONSuccess(message string) []byte {
	jsonSucc, err := json.Marshal(Success{Message: message})
	if err != nil {
		return []byte("")
	}
	return jsonSucc
}