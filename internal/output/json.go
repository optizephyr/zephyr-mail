package output

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintJSON(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		PrintError(err)
		return
	}
	fmt.Fprintln(os.Stdout, string(data))
}

func PrintError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}
