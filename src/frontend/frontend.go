package frontend

import (
	"embed"
	"io/fs"
)

var (
	//go:embed dist/*
	distEmbedFS embed.FS
	DistFS      = func() fs.FS {
		filesystem, err := fs.Sub(distEmbedFS, "dist")
		if err != nil {
			panic(err)
		}
		return filesystem
	}()
)
