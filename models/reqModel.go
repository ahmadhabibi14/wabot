package models

type Request struct {
	ModelRequest string  `json:"model"`
	Prompt       string  `json:"prompt"`
	Temperature  float64 `json:"temperature"`
	MaxTokens    int     `json:"max_tokens"`
}
