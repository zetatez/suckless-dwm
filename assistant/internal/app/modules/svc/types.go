package svc

type FormatRequest struct {
	Language string `json:"language" binding:"required"`
}

type NoteRequest struct {
	Type string `json:"type" binding:"required"`
}

type DatetimeConvertRequest struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
}

type ToggleRequest struct {
	Process string `json:"process" binding:"required"`
	Match   string `json:"match" binding:"required"`
}

type LaunchRequest struct {
	Command string `json:"command" binding:"required"`
}

type QueryRequest struct {
	Query string `json:"query" binding:"required"`
}

type OpenURLRequest struct {
	Browser string `json:"browser"`
	URL     string `json:"url" binding:"required"`
}

type DirRequest struct {
	Dir string `json:"dir"`
}
