package admin

// InstanceList represent the kind InstanceList configuration struct file
type InstanceList struct {
	Kind      string            `yaml:"kind" json:"kind"`
	Instances map[string]string `yaml:"instances" json:"instances"`
}

// Tenant represents the kind Tenant configuration struct file
type Tenant struct {
	Kind      string            `yaml:"kind" json:"kind"`
	Spec      map[string]string `yaml:"spec" json:"spec"`
	Instances []string          `yaml:"instances" json:"instances"`
}

// TenantList represents the system tenant list
type TenantList struct {
	Kind    string   `yaml:"kind" json:"kind"`
	Tenants []string `yaml:"tenants" json:"tenants"`
}
