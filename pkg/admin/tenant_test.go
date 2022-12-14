package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasMeetingPool(t *testing.T) {
	tenant := &Tenant{
		Spec: &TenantSpec{},
	}

	t.Run("tenant does not have a meeting pool", func(t *testing.T) {
		assert.False(t, tenant.HasMeetingPool())
	})

	t.Run("tenant have a meeting pool", func(t *testing.T) {
		pool := int64(10)
		tenant.Spec.MeetingsPool = &pool
		assert.True(t, tenant.HasMeetingPool())
	})
}

func TestHasUserPool(t *testing.T) {
	tenant := &Tenant{
		Spec: &TenantSpec{},
	}

	t.Run("tenant does not have a user pool", func(t *testing.T) {
		assert.False(t, tenant.HasUserPool())
	})

	t.Run("tenant have a user pool", func(t *testing.T) {
		pool := int64(100)
		tenant.Spec.UserPool = &pool
		assert.True(t, tenant.HasUserPool())
	})
}
