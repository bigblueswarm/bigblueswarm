package balancer

// InstanceStatus represents a cluster instance status at a time
type InstanceStatus struct {
	Host               string  `json:"host"`
	CPU                float64 `json:"cpu"`
	Mem                float64 `json:"mem"`
	ActiveMeeting      int64   `json:"active_meetings"`
	ActiveParticipants int64   `json:"active_participants"`
	APIStatus          string  `json:"api_status"`
}
