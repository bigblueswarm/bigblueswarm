// Package app is the bigblueswarm core
package app

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func (s *Server) clearRecordings() error {
	return s.Mapper.DeleteAll(RecodingPattern())
}

func toDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		panic(err)
	}

	return duration
}

func (s *Server) pollRecordings() {
	logger := log.WithField("context", "poll_recorder")
	logger.Info("polling recordings")
	if err := s.clearRecordings(); err != nil {
		logger.Errorln("failed to clear recordings.", err)
		return
	}

	instances, err := s.InstanceManager.ListInstances()
	if err != nil {
		logger.Errorln("failed to retrieve instances.", err)
		return
	}

	for _, instance := range instances {
		iLogger := logger.Dup().WithField("instance", instance.URL)
		recordings, err := instance.GetRecordings("")
		if err != nil {
			iLogger.Errorln("failed to retrieve recordings.", err)
			continue
		}

		for _, recording := range recordings.Recordings {
			if err := s.Mapper.Add(RecordingMapKey(recording.RecordID), instance.URL); err != nil {
				iLogger.Dup().WithField("record_id", recording).Errorln("failed to store record.", err)
				continue
			}
		}
	}
}

func (s *Server) launchRecordingPoller() {
	ticker := time.NewTicker(toDuration(s.Config.BigBlueSwarm.RecordingsPollInterval))
	for range ticker.C {
		s.pollRecordings()
	}
}
