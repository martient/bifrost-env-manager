package environmentmanager

import (
	log "github.com/sirupsen/logrus"
)

func intOrDefault(value interface{}, defaultValue int) int {
	if v, ok := value.(float64); ok {
		return int(v)
	}
	return defaultValue
}

func stringOrDefault(value interface{}, defaultValue string) string {
	if v, ok := value.(string); ok {
		return v
	}
	return defaultValue
}

func boolOrDefault(value interface{}, defaultValue bool) bool {
	if v, ok := value.(bool); ok {
		return v
	}
	return defaultValue
}

func logInfo(message string, args interface{}) {
	log.WithFields(log.Fields{
		"service": "Environment manager",
	}).Infof(message, args)
}

func logWarning(message string, args interface{}) {
	log.WithFields(log.Fields{
		"service": "Environment manager",
	}).Warnf(message, args)
}

func logError(message string, args interface{}) {
	log.WithFields(log.Fields{
		"service": "Environment manager",
	}).Errorf(message, args)
}

func logDebug(message string, args interface{}) {
	log.WithFields(log.Fields{
		"service": "Environment manager",
	}).Debugf(message, args)
}
