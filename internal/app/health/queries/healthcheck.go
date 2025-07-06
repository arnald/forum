package queries

const (
	StatusUp   string = "UP"
	StatusDown string = "DOWN"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp,omitzero"`
}
