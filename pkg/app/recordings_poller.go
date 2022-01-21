package app

import (
	"fmt"
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
	log.Info("Polling recordings")
	if err := s.clearRecordings(); err != nil {
		log.Errorln("Failed to clear recordings.", err)
		return
	}

	instances, err := s.InstanceManager.ListInstances()
	if err != nil {
		log.Errorln("Failed to retrieve instances.", err)
		return
	}

	for _, instance := range instances {
		recordings, err := instance.GetRecordings("")
		if err != nil {
			log.Errorln(fmt.Sprintf("Failed to retrieve recordings for instance %s.", instance.URL), err)
			continue
		}

		for _, recording := range recordings.Recordings {
			if err := s.Mapper.Add(RecordingMapKey(recording.RecordID), instance.URL); err != nil {
				log.Errorln(fmt.Sprintf("Failed to store recording %s.", recording.RecordID), err)
				continue
			}
		}
	}
}

func (s *Server) launchRecordingPoller() {
	ticker := time.NewTicker(toDuration(s.Config.BigBlueButton.RecordingsPollInterval))
	for range ticker.C {
		s.pollRecordings()
	}
}
