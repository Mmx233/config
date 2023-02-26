package config

import (
	"errors"
	"github.com/Mmx233/tool"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
)

var IsNewConfig = errors.New("default config generated")

type Options struct {
	Config      interface{}
	Default     interface{}
	Path        string
	FillDefault bool
	Overwrite   bool
}

type Config struct {
	opt *Options
}

func NewConfig(s *Options) *Config {
	if s.Path == "" {
		s.Path = "Config.yaml"
	}
	return &Config{opt: s}
}

func (a *Config) Init(s *Options) {
	a.opt = NewConfig(s).opt
}

func (a *Config) Save() error {
	f, e := os.OpenFile(a.opt.Path, os.O_WRONLY|os.O_CREATE, 0600)
	if e != nil {
		return e
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(a.opt.Config)
}

func (a *Config) Load() error {
	//config not exist
	if exist, e := tool.File.Exists(a.opt.Path); e != nil {

	} else if !exist {
		d, e := yaml.Marshal(a.opt.Default)
		if e != nil {
			return e
		}
		if e = os.WriteFile(a.opt.Path, d, 0600); e != nil {
			return e
		}
		return IsNewConfig
	}

	//read config
	f, e := os.OpenFile(a.opt.Path, os.O_RDONLY, 0600)
	if e != nil {
		return e
	}
	defer f.Close()
	e = yaml.NewDecoder(f).Decode(a.opt.Config)
	if e != nil {
		return e
	}

	// fill config with default value
	if a.opt.FillDefault {
		Config := reflect.ValueOf(a.opt.Config).Elem()
		Default := reflect.ValueOf(a.opt.Default).Elem()
		for i := 0; i < Default.NumField(); i++ {
			if !Default.Field(i).IsZero() && Config.Field(i).IsZero() {
				Config.Field(i).Set(Default.Field(i))
			}
		}
	}

	// fill back to ensure update
	if a.opt.Overwrite {
		return a.Save()
	}

	return nil
}

func IsNew(e error) bool {
	return e == IsNewConfig
}
