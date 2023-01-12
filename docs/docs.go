package docs

import "embed"

//go:embed static
var Docs embed.FS

//go:embed static/index.temp
var UiTemplate embed.FS
