package main

import (
	"snippetbox.org/pkg/models"

	"github.com/alexedwards/scs"
)

// Add a new StaticDir field to our application dependencies.
type App struct {
	Addr     	 string
	Database *models.Database
	HTMLDir   string
	Sessions *scs.Manager
	StaticDir string
	TLSCert   string
	TLSKey    string
}