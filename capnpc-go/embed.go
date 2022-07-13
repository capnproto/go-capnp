package main

import (
	"embed"
	"text/template"
)

var (
	//go:embed templates/*
	templateFS embed.FS

	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"title": title.String,
	}).ParseFS(templateFS, "templates/*"))
)
