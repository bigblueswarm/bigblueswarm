// Package app is the bigblueswarm core
package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// RequestLogger is a custom request logger initialized with a request identifier
type RequestLogger struct {
	*log.Entry
}

func setLogger(c *gin.Context) {
	c.Set("logger", newRequestLogger())
}

func getLogger(c *gin.Context) *RequestLogger {
	return c.MustGet("logger").(*RequestLogger)
}

func newRequestLogger() *RequestLogger {
	return &RequestLogger{
		log.WithFields(log.Fields{
			"request_id": uuid.New().String(),
		}),
	}
}

func (rl *RequestLogger) setFields(fields log.Fields) *RequestLogger {
	for k, v := range fields {
		rl.Data[k] = v
	}

	return rl
}

func (rl *RequestLogger) addField(key string, value string) *RequestLogger {
	rl.Data[key] = value
	return rl
}

func (rl *RequestLogger) dup() *RequestLogger {
	return &RequestLogger{
		rl.Dup(),
	}
}
