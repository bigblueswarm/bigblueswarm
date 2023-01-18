// Package balancer manage the balancer progress and choose the next server
package balancer

// InstanceStatus represents a cluster instance status at a time
type InstanceStatus struct {
	Host         string  `json:"host"`
	CPU          float64 `json:"cpu"`
	Mem          float64 `json:"mem"`
	Meetings     int64   `json:"meetings"`
	Participants int64   `json:"participants"`
	APIStatus    string  `json:"api_status"`
}
