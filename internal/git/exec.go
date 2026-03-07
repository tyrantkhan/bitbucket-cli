package git

import (
	"fmt"
	"os"
	"os/exec"
)

// Clone runs git clone with the given URL and optional directory.
func Clone(url string, dir string, args ...string) error {
	cmdArgs := []string{"clone", url}
	if dir != "" {
		cmdArgs = append(cmdArgs, dir)
	}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("git", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Cloning into '%s'...\n", url)
	return cmd.Run()
}
