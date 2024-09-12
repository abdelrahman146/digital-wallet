package errs

type CustomError struct {
	Message  string
	Code     string            `json:"code"`
	HttpCode int               `json:"httpCode"`
	Original error             `json:"error,omitempty"`
	Fields   map[string]string `json:"fields,omitempty"`
}

func (e CustomError) Error() string {
	return e.Message
}
