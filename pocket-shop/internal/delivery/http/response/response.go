package response

type ErrorResponse struct {
	Error   string   `json:"error"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}
