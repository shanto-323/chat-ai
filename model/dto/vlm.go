package dto


type VLMResult struct {
	ImageURL   string   `json:"image_url"`
	Summary    string   `json:"summary"`
	Objects    []string `json:"objects"`
	Confidence float32  `json:"confidence"`
}
