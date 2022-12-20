// Package app is the bigblueswarm core
package app

import "github.com/bigblueswarm/bigblueswarm/v2/pkg/api"

func noInstanceFoundError() *api.Error {
	return api.CreateError("noInstanceFound", "BigBlueSwarm do not find a valid BigBlueButton instance for your request")
}

func tenantNotFoundError() *api.Error {
	return api.CreateError("tenantNotFound", "BigBlueSwarm does not find the requesting tenant")
}

func serverError(message string) *api.Error {
	return api.CreateError("internalError", message)
}

func getTenantError() *api.Error {
	return serverError("BigBlueSwarm failed to retrieve the requesting tenant")
}

func meetingPoolReachedError() *api.Error {
	return api.CreateError("meetingPoolReached", "Your tenant reached the meeting pool limit and can't create a new one.")
}

func userPoolReachedError() *api.Error {
	return api.CreateError("userPoolReached", "Your tenant reached the user pool limit.")
}
