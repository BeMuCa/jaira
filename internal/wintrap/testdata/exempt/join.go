package fixture

import "strings"

func branch(b, remote string) bool {
	// A git refname is slash-separated on every platform.
	//wintrap:ok
	return strings.HasPrefix(b, remote+"/")
}
