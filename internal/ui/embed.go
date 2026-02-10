package ui

import (
	"embed"
	"io/fs"
)

//go:embed static/*
var embeddedFiles embed.FS

func GetFileSystem() (fs.FS, error) {
	return fs.Sub(embeddedFiles, "static")
}
