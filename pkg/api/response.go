package api

import "digital-wallet/pkg/errs"

func NewErrorResponse(err error) (httpCode int, resp ErrorResponse) {
	customErr := errs.HandleError(err)
	httpCode = customErr.HttpCode
	body := ErrorResponseBody{
		Message:  customErr.Message,
		HttpCode: customErr.HttpCode,
		Code:     customErr.Code,
	}
	if customErr.Original != nil {
		body.Reason = customErr.Original.Error()
	}
	if customErr.Fields != nil {
		body.Fields = customErr.Fields
	}
	resp = ErrorResponse{
		Success: false,
		Error:   body,
	}
	return httpCode, resp
}

func NewSuccessResponse(result interface{}) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Result:  result,
	}
}
