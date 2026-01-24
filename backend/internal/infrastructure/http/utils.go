package http

import (
	"errors"
	"net/http"

	"github.com/Haya372/web-app-template/backend/internal/domain/vo"
)

type errorResponse struct {
	Code    vo.ErrorCode `json:"code"`
	Message string       `json:"message"`
	Details any          `json:"details"`
}

func handleError(err error) (int, any) {
	var baseErr vo.Error
	if errors.As(err, &baseErr) {
		return baseErr.Status(), errorResponse{
			Code:    baseErr.Code(),
			Message: baseErr.Message(),
			Details: baseErr.Details(),
		}
	}

	return http.StatusInternalServerError, errorResponse{
		Code:    vo.InternalErrorCode,
		Message: "internal server error",
	}
}
