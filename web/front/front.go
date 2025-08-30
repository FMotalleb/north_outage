package front

import (
	"embed"
	"io/fs"
)

//go:generate npm i
//go:generate npm run build

//go:embed dist/*
var distFS embed.FS

// GetDist returns a filesystem rooted at the embedded "dist" directory.
// It enables access to static files embedded at compile time.
func GetDist() (fs.FS, error) {
	return fs.Sub(distFS, "dist")
}
