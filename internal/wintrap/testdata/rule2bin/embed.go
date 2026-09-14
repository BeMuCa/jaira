package fixture

import "embed"

//go:embed assets/*.png
var assetFS embed.FS
