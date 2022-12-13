// Package admin manages the bigblueswarm admin part
package admin

// InstanceList represent the kind InstanceList configuration struct file
type InstanceList struct {
	Kind      string            `yaml:"kind" json:"kind"`
	Instances map[string]string `yaml:"instances" json:"instances"`
}

// TenantSpec represents the tenant spec configuration
type TenantSpec struct {
	Host         string `yaml:"host,omitempty" json:"host,omitempty"`
	Secret       string `yaml:"secret,omitempty" json:"secret,omitempty"`
	MeetingsPool *int64 `yaml:"meetings_pool,omitempty" json:"meetings_pool,omitempty"`
	UserPool     *int64 `yaml:"user_pool,omitempty" json:"user_pool,omitempty"`
}

// Tenant represents the kind Tenant configuration struct file
type Tenant struct {
	Kind      string      `yaml:"kind" json:"kind"`
	Spec      *TenantSpec `yaml:"spec" json:"spec"`
	Instances []string    `yaml:"instances" json:"instances"`
}

// TenantList represents the system tenant list
type TenantList struct {
	Kind    string             `yaml:"kind" json:"kind"`
	Tenants []TenantListObject `yaml:"tenants" json:"tenants"`
}

// TenantListObject represents a Tenant in a TenantList
type TenantListObject struct {
	Hostname      string `json:"hostname"`
	InstanceCount int    `json:"instance_count"`
}
