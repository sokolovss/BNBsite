package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/sokolovss/BNBsite/internal/models"
	"html/template"
	"log"
)

//AppConfig holds the app config
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	IsProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
