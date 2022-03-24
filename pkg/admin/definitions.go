package admin

// InstanceList represent the kind InstanceList configuration struct file
type InstanceList struct {
	Kind      string            `yaml:"kind"`
	Instances map[string]string `yaml:"instances"`
}
