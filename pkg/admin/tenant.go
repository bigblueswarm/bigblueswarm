// Package admin manages the bigblueswarm admin part
package admin

// HasMeetingPool check if tenant as a meeting pool constraint
func (t *Tenant) HasMeetingPool() bool {
	return t.Spec.MeetingsPool != nil
}

// HasUserPool check if tenant as a user pool constraint
func (t *Tenant) HasUserPool() bool {
	return t.Spec.UserPool != nil
}
