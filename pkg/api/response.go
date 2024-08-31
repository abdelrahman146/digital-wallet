package api

import "digital-wallet/pkg/errs"

func NewErrorResponse(err error) (httpCode int, resp ErrorResponse) {
	customErr := errs.HandleError(err)
	httpCode = customErr.HttpCode
	resp = ErrorResponse{
		Success: false,
		Error: ErrorResponseBody{
			Message:  customErr.Desc,
			HttpCode: customErr.HttpCode,
			Code:     customErr.Code,
			Reason:   customErr.Original.Error(),
		},
	}
	return httpCode, resp
}

func NewSuccessResponse(result interface{}) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Result:  result,
	}
}
