package images

type ImageFactResponse struct {
	Services  []string `json:"services"`
	Queues    []string `json:"queues"`
	Databases []string `json:"databases"`
}
