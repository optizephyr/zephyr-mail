package output

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintJSON(v any) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		PrintError(err)
	}
}

func PrintError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}
