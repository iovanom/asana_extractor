package models

type Project struct {
	ID          string `json:"gid"`
	Name        string `json:"name"`
	Archived    bool   `json:"archived"`
	Completed   bool   `json:"completed"`
	CompletedBy string `json:"completed_by"`
	CompletedAt string `json:"completed_at"`
}
