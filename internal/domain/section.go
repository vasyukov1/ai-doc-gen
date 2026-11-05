package domain

type DocumentTable struct {
	Caption string     `json:"caption"`
	Columns []string   `json:"columns"`
	Rows    [][]string `json:"rows"`
}

type Section struct {
	Tables    []DocumentTable `json:"tables,omitempty"`
	Origin    string          `json:"origin,omitempty"`
	Provider  string          `json:"provider,omitempty"`
	Model     string          `json:"model,omitempty"`
	ID        string          `json:"id"`
	Heading   string          `json:"heading"`
	Body      string          `json:"body"`
	Citations []string        `json:"citations"`
	Missing   []string        `json:"missing_information"`
	Edited    bool            `json:"edited"`
	CacheKey  string          `json:"cache_key"`
}

type Usage struct {
	InputTokens  int `json:"prompt_tokens"`
	OutputTokens int `json:"completion_tokens"`
}
