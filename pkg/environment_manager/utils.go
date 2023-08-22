package environmentmanager

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
