package output

import (
	"charm.land/lipgloss/v2"
)

var (
	// BitbucketBlue is the standard Bitbucket brand color.
	BitbucketBlue = lipgloss.Color("#0052CC")
	// Bitbucket blue
	Blue = lipgloss.NewStyle().Foreground(BitbucketBlue)
	// Status colors
	Green   = lipgloss.NewStyle().Foreground(lipgloss.Color("#36B37E"))
	Red     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5630"))
	Yellow  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAB00"))
	Cyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00B8D9"))
	Muted   = lipgloss.NewStyle().Foreground(lipgloss.Color("#97A0AF"))
	Bold    = lipgloss.NewStyle().Bold(true)
	Header  = lipgloss.NewStyle().Bold(true).Foreground(BitbucketBlue)
	Success = Green
	Error   = Red
	Warning = Yellow
)

// StatusColor returns the appropriate style for a given status string.
func StatusColor(status string) lipgloss.Style {
	switch status {
	case "OPEN", "SUCCESSFUL", "PASSED", "approved":
		return Green
	case "MERGED", "COMPLETED":
		return Cyan
	case "DECLINED", "FAILED", "ERROR", "STOPPED":
		return Red
	case "PENDING", "PAUSED", "HALTED":
		return Yellow
	case "IN_PROGRESS", "RUNNING":
		return Blue
	default:
		return Muted
	}
}
