package main

import (
	"embed"
	"strings"
	"text/template"
)

var (
	//go:embed templates/*
	templateFS embed.FS

	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"title": strings.Title,
	}).ParseFS(templateFS, "templates/*"))
)
