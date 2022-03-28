package admin

// InstanceList represent the kind InstanceList configuration struct file
type InstanceList struct {
	Kind      string            `yaml:"kind"`
	Instances map[string]string `yaml:"instances"`
}

// Tenant represents the kind Tenant configuration struct file
type Tenant struct {
	Kind      string            `yaml:"kind"`
	Spec      map[string]string `yaml:"spec"`
	Instances []string          `yaml:"instances"`
}
