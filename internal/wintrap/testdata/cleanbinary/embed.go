package fixture

import "embed"

//go:embed assets/*.png docs/*.txt
var assetFS embed.FS
