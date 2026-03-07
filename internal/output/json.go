package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// RenderJSON outputs the given value as indented JSON.
func RenderJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(data))
	return nil
}
