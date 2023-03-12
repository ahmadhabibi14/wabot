package models

type TextCompletionResponse struct {
	ID            string    `json:"id"`
	Object        string    `json:"object"`
	Created       int       `json:"created"`
	ModelResponse int       `json:"model"`
	Choices       []Choices `json:"choices"`
	Usage         Usage     `json:"usage"`
}

type Choices struct {
	Text          string `json:"text"`
	Index         int    `json:"index"`
	Logprobs      string `json:"logprobs"`
	Finish_reason string `json:"finish_reason"`
}

type Usage struct {
	PrompToken       int `json:"prompt_token"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
