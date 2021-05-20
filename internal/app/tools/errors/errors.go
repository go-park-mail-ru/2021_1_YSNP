package errors

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/status"
)

type ErrorType uint8

const (
	InternalError ErrorType = iota + 1
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
	ProductAlreadyLiked
	PromoteEmptyLabel
	InvalidCSRFToken
	EmptySearch
	ProductClose
	WrongOwner
	ChatNotExist
	NoPhoto
	ReviewExist
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
		HttpError: http.StatusUnauthorized,
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
	ProductAlreadyLiked: {
		ErrorCode: ProductAlreadyLiked,
		HttpError: http.StatusBadRequest,
		Message:   "product already liked",
	},
	EmptySearch: {
		ErrorCode: EmptySearch,
		HttpError: http.StatusNotFound,
		Message:   "searching products dont't exist",
	},
	PromoteEmptyLabel: {
		ErrorCode: PromoteEmptyLabel,
		HttpError: http.StatusBadRequest,
		Message:   "promote label doesn't exist",
	},
	InvalidCSRFToken: {
		ErrorCode: InvalidCSRFToken,
		HttpError: http.StatusForbidden,
		Message:   "Forbidden - CSRF token invalid",
	},
	ProductClose: {
		ErrorCode: ProductClose,
		HttpError: http.StatusBadRequest,
		Message:   "Product already close",
	},
	WrongOwner: {
		ErrorCode: WrongOwner,
		HttpError: http.StatusForbidden,
		Message:   "Forbidden - wrong owner",
	},
	ChatNotExist: {
		ErrorCode: ChatNotExist,
		HttpError: http.StatusNotFound,
		Message:   "chat doesn't exist",
	},
	NoPhoto: {
		ErrorCode: NoPhoto,
		HttpError: http.StatusBadRequest,
		Message:   "no photo",
	},
	ReviewExist: {
		ErrorCode: ReviewExist,
		HttpError: http.StatusBadRequest,
		Message:   "user has already left a review on this product",
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

func GRPCError(err error) *Error {
	grpcErr, has := CustomErrors[ErrorType(status.Code(err))]
	grpcErr.Message = status.Convert(err).Message()
	if !has {
		// for grpc connection error
		return UnexpectedInternal(err)
	}
	return grpcErr
}
