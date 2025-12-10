package model

type VLMResult struct {
	ImageIndex int       `json:"image_index"`
	Summary    string    `json:"summary"`
	Objects    []string  `json:"objects"`
	SceneDesc  string    `json:"scene_description"`
	Embeddings []float32 `json:"embeddings,omitempty"`
	Confidence float32   `json:"confidence"`
}
