package config

import (
	"errors"
	"github.com/Mmx233/tool"
	"gopkg.in/yaml.v3"
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

func Load(s *Options) error {
	if s.Path == "" {
		s.Path = "Config.yaml"
	}

	//config not exist
	if !tool.File.Exists(s.Path) {
		d, e := yaml.Marshal(s.Default)
		if e != nil {
			return e
		}
		if e = tool.File.Write(s.Path, d); e != nil {
			return e
		}
		return NewConfig
	}

	//read config
	d, e := tool.File.Read(s.Path)
	if e != nil {
		return e
	}
	e = yaml.Unmarshal(d, s.Config)
	if e != nil {
		return e
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
		d, e = yaml.Marshal(s.Config)
		if e != nil {
			return e
		}
		if e = tool.File.Write(s.Path, d); e != nil {
			return e
		}
	}

	return nil
}

func IsNew(e error) bool {
	return e == NewConfig
}
