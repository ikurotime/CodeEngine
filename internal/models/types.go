package models

type ErrorResponse struct {
	Error string `json:"error"`
}

type ExecuteRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type ExecuteResponse struct {
	Output string `json:"output"`
}

var LanguageToExtension = map[string]string{
	"python3": "py",
}
