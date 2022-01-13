package api

import "time"

type FileItem struct {
	EstimatedRowsCount int    `json:"estimated_rows_count"`
	Filename           string `json:"filename"`
}

type StatusItem struct {
	State               string  `json:"state"`
	Progress            float32 `json:"progress"`
	Message             string  `json:"message"`
	TimeLeftDescription string  `json:"time_left_description,omitempty"`
	TimeLeftSeconds     int     `json:"time_left_seconds,omitempty"`
}

type SpreadsheetJob struct {
	Id          int        `json:"id"`
	Fields      []string   `json:"fields"`
	File        FileItem   `json:"file"`
	Status      StatusItem `json:"status"`
	DownloadUrl string     `json:"download_url"`
	ExpiresAt   time.Time  `json:"expires_at"`
	Error       string     `json:"error,omitempty"`
}
