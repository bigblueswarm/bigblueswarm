package test

// Cluster represents the test cluster for the application.
type Cluster struct {
	InfluxDB       *Container
	Redis          *Container
	BigBlueButtons []*Container
}
