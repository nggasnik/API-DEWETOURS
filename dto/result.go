package dto

// membuat struct untuk response ketika request berhasil diproses
type SuccessResult struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// membuat struct untuk response ketika terjadi error
type ErrorResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
