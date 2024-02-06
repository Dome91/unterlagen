package views

import "embed"

//go:embed public/*
var Assets embed.FS

//go:embed templates/*
var Templates embed.FS
