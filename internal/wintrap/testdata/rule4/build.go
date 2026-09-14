package fixture

import (
	"os/exec"
	"path/filepath"
)

func build(dir string) error {
	bin := filepath.Join(dir, "jaira")
	return exec.Command("go", "build", "-o", bin, "./cmd/jaira").Run()
}
