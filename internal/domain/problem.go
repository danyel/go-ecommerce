package domain

type ProblemDetail struct {
	Type     string         `json:"type,omitempty"`
	Title    string         `json:"title"`
	Status   int            `json:"status"`
	Details  string         `json:"detail,omitempty"`
	Instance string         `json:"instance,omitempty"`
	Errors   map[string]any `json:"errors,omitempty"`
}

func (problemDetail *ProblemDetail) IsSuccess() bool {
	return problemDetail.Status > 199 && problemDetail.Status < 300
}

func (problemDetail *ProblemDetail) IsError() bool {
	return problemDetail.Status > 399
}

func (problemDetail *ProblemDetail) Error() string {
	return problemDetail.Details
}
