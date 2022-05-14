package internal

import (
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	Database struct {
		URL  string `envconfig:"DATABASE_URL" default:"postgres://postgres:postgrespw@postgres:5432"`
		PORT string `envconfig:"PORT" default:":7008"`
	}
	// here could be some other configs
}

func NewSettings() (*Settings, error) {
	var settings Settings
	err := envconfig.Process("", &settings)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}
