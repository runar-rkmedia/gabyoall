package frontend

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
)

//go:embed dist/*
var DistFiles embed.FS

type myFS struct {
	content embed.FS
}

func (c myFS) Open(name string) (fs.File, error) {
	return c.content.Open(path.Join("dist", name))
}

var DistServer = http.FileServer(http.FS(myFS{DistFiles}))
