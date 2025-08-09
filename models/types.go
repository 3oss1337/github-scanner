package models

type FileInfo struct {
	Path string `json:"path"`
	Type string `json:"type"`
}
type Pattern struct {
	Name  string
	Regex string
}
type ScanResult struct {
	Type  string
	Match string
}
type ScanRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}
type FileContent struct {
	Content string `json:"content"`
}
type ScanResponse struct {
	File  string `json:"file"`
	Type  string `json:"type"`
	Match string `json:"match"`
}
