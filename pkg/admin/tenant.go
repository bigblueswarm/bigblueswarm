// Package admin manages the bigblueswarm admin part
package admin

func (t *Tenant) HasMeetingPool() bool {
	return t.Spec.MeetingsPool != nil
}

func (t *Tenant) HasUserPool() bool {
	return t.Spec.UserPool != nil
}
