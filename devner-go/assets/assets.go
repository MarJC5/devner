package assets

import "embed"

//go:embed compose/*.yaml dockerfiles/*
var FS embed.FS
