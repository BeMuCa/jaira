package cli

import (
	"fmt"
	"strings"
)

func announceLine(written []string) string {
	if len(written) == 0 {
		return ""
	}
	return fmt.Sprintf("Told your agents about this board: %s\n", strings.Join(written, ", "))
}
