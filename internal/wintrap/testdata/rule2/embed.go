package fixture

import "embed"

//go:embed prompts/*.md
var promptFS embed.FS
