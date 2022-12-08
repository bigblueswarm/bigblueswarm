// Package utils provide few utilies functions
package utils

import (
	"fmt"
	"strings"
)

// FormatInstancesFilter format influxdb filter for bigblueswarm_host from an instance
// list like r["bigblueswarm_host"] == "http://localhost/bigbluebutton" or r["bigblueswarm_host"] == "http://localhost:8080/bigbluebutton"
func FormatInstancesFilter(instances []string) string {
	var result string
	for i, instance := range instances {
		filter := fmt.Sprintf(`r["bigblueswarm_host"] == "%s"`, instance)
		result = fmt.Sprintf("%s %s", result, filter)

		if i != (len(instances) - 1) {
			result = fmt.Sprintf("%s or", result)
		}
	}

	return strings.TrimSpace(result)
}
