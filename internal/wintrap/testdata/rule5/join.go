package fixture

import (
	"os"
	"strings"
)

func setup(root string) error {
	return os.MkdirAll(root+"/.jaira", 0o755)
}

func label(root, p string) string {
	return strings.TrimPrefix(p, root+"/")
}
