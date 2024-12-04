package web

import (
	"embed"
	"html/template"
)

//go:embed templates/*
var views embed.FS

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("").Funcs(template.FuncMap{
		"DerefBool": func(ptr *bool) bool { return *ptr },
	}).ParseFS(views, "templates/*.html"))
}

func GetTemplate() *template.Template {
	return tmpl
}
