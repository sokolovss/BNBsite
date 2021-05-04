package config

import "html/template"

//AppConfig holds the app config
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
}
