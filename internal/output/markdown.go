package output

import (
	"fmt"

	"charm.land/glamour/v2"
)

// RenderMarkdown renders a markdown string for terminal display.
func RenderMarkdown(md string) {
	if md == "" {
		return
	}

	rendered, err := glamour.Render(md, "dark")
	if err != nil {
		// Fallback to raw output
		fmt.Println(md)
		return
	}
	fmt.Print(rendered)
}
