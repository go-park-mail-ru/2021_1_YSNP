package errors

import (
	"encoding/json"
	"net/http"
)

type ErrorType uint8

const (
	InternalError ErrorType = iota
	BadRequest
	UserNotExist
	WrongPassword
	TelephoneAlreadyExists
	SessionNotExist
	SessionExpired
	WrongErrorCode
	UserUnauthorized
	EmptyContext
	ProductNotExist
)

type Error struct {
	ErrorCode ErrorType `json:"code"`
	HttpError int       `json:"-"`
	Message   string    `json:"message"`
}

type Success struct {
	Message string `json:"message"`
	ID      uint64 `json:"id,omitempty"`
}

//func (e Error) Error() string {
//	return e.Message
//}

func JSONError(error *Error) []byte {
	jsonError, err := json.Marshal(error)
	if err != nil {
		return []byte("")
	}
	return jsonError
}

func JSONSuccess(message ...interface{}) []byte {
	if len(message) > 1 {
		jsonSucc, err := json.Marshal(Success{Message: message[0].(string), ID: message[1].(uint64)})
		if err != nil {
			return []byte("")
		}
		return jsonSucc
	}
	jsonSucc, err := json.Marshal(Success{Message: message[0].(string)})
	if err != nil {
		return []byte("")
	}
	return jsonSucc
}

var CustomErrors = map[ErrorType]*Error{
	InternalError: {
		ErrorCode: InternalError,
		HttpError: http.StatusInternalServerError,
		Message:   "somthing wrong",
	},
	BadRequest: {
		ErrorCode: BadRequest,
		HttpError: http.StatusBadRequest,
		Message:   "wrong request",
	},
	UserNotExist: {
		ErrorCode: UserNotExist,
		HttpError: http.StatusBadRequest,
		Message:   "user with this telephone doesn't exist",
	},
	WrongPassword: {
		ErrorCode: WrongPassword,
		HttpError: http.StatusBadRequest,
		Message:   "wrong password",
	},
	TelephoneAlreadyExists: {
		ErrorCode: TelephoneAlreadyExists,
		HttpError: http.StatusBadRequest,
		Message:   "user with this telephone already exists",
	},
	SessionNotExist: {
		ErrorCode: SessionNotExist,
		HttpError: http.StatusBadRequest,
		Message:   "user session doesn't exists",
	},
	SessionExpired: {
		ErrorCode: SessionExpired,
		HttpError: http.StatusUnauthorized,
		Message:   "user session expired",
	},
	WrongErrorCode: {
		ErrorCode: WrongErrorCode,
		HttpError: http.StatusInternalServerError,
		Message:   "can't specify error",
	},
	UserUnauthorized: {
		ErrorCode: UserUnauthorized,
		HttpError: http.StatusUnauthorized,
		Message:   "user unauthorized",
	},
	EmptyContext: {
		ErrorCode: EmptyContext,
		HttpError: http.StatusInternalServerError,
		Message:   "empty context",
	},
	ProductNotExist: {
		ErrorCode: ProductNotExist,
		HttpError: http.StatusNotFound,
		Message:   "product doesn't exist",
	},
}

func Cause(code ErrorType) *Error {
	err, ok := CustomErrors[code]
	if !ok {
		return CustomErrors[WrongErrorCode]
	}
	return err
}

func UnexpectedInternal(err error) *Error {
	unexpErr := CustomErrors[InternalError]
	unexpErr.Message = err.Error()

	return unexpErr
}

func UnexpectedBadRequest(err error) *Error {
	unexpErr := CustomErrors[BadRequest]
	unexpErr.Message = err.Error()

	return unexpErr
}
