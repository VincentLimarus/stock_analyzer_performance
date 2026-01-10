package model

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
	StatusError   = "error"
)

// Response using jsend format
// ref: https://github.com/omniti-labs/jsend
type Response struct {
	Status       any    `json:"status"`
	Message      string `json:"message"`
	Code         string `json:"code,omitempty"`
	Data         any    `json:"data,omitempty"`
	Limit        uint   `json:"limit,omitempty"`
	TotalRecords uint64 `json:"totalRecords,omitempty"`
	CurrentPage  uint   `json:"currentPage,omitempty"`
	NextPage     uint   `json:"nextPage,omitempty"`
	PreviousPage uint   `json:"previousPage,omitempty"`
	TotalPages   uint   `json:"totalPages,omitempty"`
}

type ValidationResponse struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}
