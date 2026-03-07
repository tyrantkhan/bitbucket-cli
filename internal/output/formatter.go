package output

// Format dispatches output to the appropriate renderer based on format.
func Format(format string, jsonData interface{}, headers []string, rows [][]string) error {
	switch format {
	case "json":
		return RenderJSON(jsonData)
	default:
		RenderTable(headers, rows)
		return nil
	}
}
