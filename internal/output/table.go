package output

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// RenderTable outputs a formatted table with headers and rows.
func RenderTable(headers []string, rows [][]string) {
	if len(rows) == 0 {
		fmt.Println(Muted.Render("No results found."))
		return
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0052CC"))
	cellStyle := lipgloss.NewStyle()

	// Render header
	var headerParts []string
	for i, h := range headers {
		headerParts = append(headerParts, headerStyle.Render(pad(h, widths[i])))
	}
	fmt.Println(strings.Join(headerParts, "  "))

	// Separator
	var sepParts []string
	for _, w := range widths {
		sepParts = append(sepParts, strings.Repeat("─", w))
	}
	fmt.Println(Muted.Render(strings.Join(sepParts, "──")))

	// Render rows
	for _, row := range rows {
		var parts []string
		for i, cell := range row {
			if i < len(widths) {
				parts = append(parts, cellStyle.Render(pad(cell, widths[i])))
			}
		}
		fmt.Println(strings.Join(parts, "  "))
	}
}

func pad(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
