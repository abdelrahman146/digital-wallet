package api

type List[T any] struct {
	Items []T   `json:"items"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type ErrorResponseBody struct {
	Message  string            `json:"message"`
	Desc     string            `json:"desc"`
	HttpCode int               `json:"httpCode"`
	Code     string            `json:"code"`
	Reason   string            `json:"reason,omitempty"`
	Fields   map[string]string `json:"fields,omitempty"`
}

type ErrorResponse struct {
	Success bool              `json:"success"`
	Error   ErrorResponseBody `json:"error"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
}
