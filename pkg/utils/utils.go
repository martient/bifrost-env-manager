package utils

import (
	log "github.com/sirupsen/logrus"
)

func IntOrDefault(value interface{}, defaultValue int) int {
	if v, ok := value.(float64); ok {
		return int(v)
	}
	return defaultValue
}

func StringOrDefault(value interface{}, defaultValue string) string {
	if v, ok := value.(string); ok {
		return v
	}
	return defaultValue
}

func BoolOrDefault(value interface{}, defaultValue bool) bool {
	if v, ok := value.(bool); ok {
		return v
	}
	return defaultValue
}

func LogInfo(message string, args interface{}, service string) {
	log.WithFields(log.Fields{
		"service": service,
	}).Infof(message, args)
}

func LogWarning(message string, args interface{}, service string) {
	log.WithFields(log.Fields{
		"service": service,
	}).Warnf(message, args)
}

func LogError(message string, args interface{}, service string) {
	log.WithFields(log.Fields{
		"service": service,
	}).Errorf(message, args)
}

func LogDebug(message string, args interface{}, service string) {
	log.WithFields(log.Fields{
		"service": service,
	}).Debugf(message, args)
}
