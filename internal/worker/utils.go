package worker

func initInt(source *int, defaultValue int) int {
	if source == nil || *source == 0 {
		return defaultValue
	}
	return *source
}

func initString(source *string, defaultValue string) string {
	if source == nil || *source == "" {
		return defaultValue
	}
	return *source
}

func initStringSlice(source *[]string, defaultValue []string) []string {
	if source == nil || len(*source) == 0 {
		return defaultValue
	}
	return *source
}
