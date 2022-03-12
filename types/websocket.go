package types

type WebsocketContact struct {
	Type   string                    `json:"type"`
	Error  *WebsocketError           `json:"error"`
	Update *WebsocketUpdateJobStatus `json:"update_job_status"`
	NewJob *Job                      `json:"new_job"`
}

type WebsocketError struct {
	Error string `json:"errpr"`
	Job   *Job   `job:"job"`
}

type WebsocketUpdateJobStatus struct {
	Job Job `json:"job"`
}
