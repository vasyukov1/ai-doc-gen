package domain

type Section struct {
	ID        string   `json:"id"`
	Heading   string   `json:"heading"`
	Body      string   `json:"body"`
	Citations []string `json:"citations"`
	Missing   []string `json:"missing_information"`
	Edited    bool     `json:"edited"`
	CacheKey  string   `json:"cache_key"`
}

type Usage struct {
	InputTokens  int `json:"prompt_tokens"`
	OutputTokens int `json:"completion_tokens"`
}
