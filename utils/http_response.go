package utils

import (
	"fmt"
	"template-go-api/enums"
	"time"

	"github.com/gin-gonic/gin"
)

type BaseHttpResponse struct {
	Code       enums.ResponseCode `json:"code"`
	Message    string             `json:"message"`
	StatusCode int                `json:"status_code"`
	Timestamp  time.Time          `json:"timestamp"`
	Payload    *Payload           `json:"payload"`
}
type HTTPResponse struct {
	BaseHttpResponse
}

type HTTPError struct {
	BaseHttpResponse
	Errors *map[string][]string `json:"errors"`
}

func (H HTTPError) Error() string {
	return fmt.Sprintf(
		"message : %s, code: %s, timestamp: %s, status_code: %d",
		H.Message,
		H.Code,
		H.Timestamp,
		H.StatusCode,
	)
}

type Payload struct {
	Data interface{}     `json:"data"`
	Meta *PaginationMeta `json:"meta,omitempty"`
}

type PaginationMeta struct {
	CurrentPage int    `json:"current_page"`
	From        int    `json:"from"`
	To          int    `json:"to"`
	LastPage    int    `json:"last_page"`
	Path        string `json:"path"`
	PerPage     int    `json:"per_page"`
	Total       int64  `json:"total"`
}

func NewHttpSuccess(message string, payload *Payload, responseCode ...enums.ResponseCode) *HTTPResponse {
	responseStatus := enums.SUCCESS
	if len(responseCode) > 0 {
		responseStatus = responseCode[0]
	}
	return &HTTPResponse{
		BaseHttpResponse: BaseHttpResponse{
			Code:       responseStatus,
			Timestamp:  time.Now(),
			StatusCode: responseStatus.HTTPStatus(),
			Message:    message,
			Payload:    payload,
		},
	}
}

func NewHttpError(message string, code enums.ResponseCode, errors *map[string][]string) *HTTPError {
	return &HTTPError{
		BaseHttpResponse: BaseHttpResponse{
			Code:       code,
			StatusCode: code.HTTPStatus(),
			Timestamp:  time.Now(),
			Message:    message,
			Payload:    nil,
		},
		Errors: errors,
	}
}

func ResponseJSON(c *gin.Context, code enums.ResponseCode, message string, data interface{}) {
	var finalPayload *Payload

	switch v := data.(type) {
	case *Payload:
		finalPayload = v
	case Payload:
		finalPayload = &v
	default:
		finalPayload = &Payload{
			Data: data,
			Meta: nil,
		}
	}

	c.JSON(
		code.HTTPStatus(),
		NewHttpSuccess(message, finalPayload, code),
	)
}
