package fixture

import "os"

func executable(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode().Perm()&0o111 != 0
}
