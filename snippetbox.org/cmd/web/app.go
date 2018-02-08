package main

import (
	"snippetbox.org/pkg/models"

	"github.com/alexedwards/scs"
)

type App struct {
	Addr      string
	Database *models.Database
	HTMLDir   string
	Sessions *scs.Manager
	Admin    *scs.Manager
	StaticDir string
	TLSCert   string
	TLSKey    string
}