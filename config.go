package config

import (
	"errors"
	"github.com/Mmx233/tool"
	"reflect"
)

var NewConfig = errors.New("default config generated")

type Options struct {
	Config      interface{}
	Default     interface{}
	Path        string
	FillDefault bool
	Overwrite   bool
}

func Load(s Options) error {
	if s.Path == "" {
		s.Path = "Config.json"
	}

	//config not exist
	if !tool.File.Exists(s.Path) {
		if err := tool.File.WriteJson(
			s.Path,
			s.Default,
		); err != nil {
			return err
		}
		return NewConfig
	}

	if err := tool.File.ReadJson(s.Path, s.Config); err != nil {
		return err
	}

	// fill config with default value
	if s.FillDefault {
		Config := reflect.ValueOf(s.Config).Elem()
		Default := reflect.ValueOf(s.Default).Elem()
		for i := 0; i < Default.NumField(); i++ {
			if !Default.Field(i).IsZero() && Config.Field(i).IsZero() {
				Config.Field(i).Set(Default.Field(i))
			}
		}
	}

	// fill back to ensure update
	if s.Overwrite {
		_ = tool.File.WriteJson(s.Path, s.Config)
	}

	return nil
}

func IsNew(e error) bool {
	return e == NewConfig
}
