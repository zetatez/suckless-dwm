package response

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess       = 0
	CodeInvalidParams = 10001
	CodeDatabaseError = 10002
	CodeThirdPartyErr = 10003
	CodeNotFound      = 10004
	CodeUnauthorized  = 10005
	CodeForbidden     = 10006
	CodeServerError   = 10000
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func Err(c *gin.Context, code int, userMsg string) {
	ErrWithInternal(c, code, userMsg, nil)
}

func ErrWithInternal(c *gin.Context, code int, userMsg string, internalErr error) {
	if internalErr != nil {
		log.Printf("[ERROR] code=%d internal=%v", code, internalErr)
	}
	status := statusFromCode(code)
	c.JSON(status, Response{
		Code:    code,
		Message: userMsg,
	})
}

func Errf(c *gin.Context, code int, format string, args ...any) {
	userMsg := fmt.Sprintf(format, args...)
	Err(c, code, userMsg)
}

func statusFromCode(code int) int {
	switch code {
	case CodeSuccess:
		return http.StatusOK
	case CodeInvalidParams:
		return http.StatusBadRequest
	case CodeDatabaseError:
		return http.StatusInternalServerError
	case CodeThirdPartyErr:
		return http.StatusBadGateway
	case CodeNotFound:
		return http.StatusNotFound
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeServerError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
