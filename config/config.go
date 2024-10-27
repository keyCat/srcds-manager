package config

import (
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"log"
	"os"
	"path/filepath"
)

type Root struct {
	Project  string   `yaml:"project"`
	Path     string   `yaml:"path"`
	Ip       string   `yaml:"ip"`
	Params   string   `yaml:"params"`
	Map      string   `yaml:"map"`
	Niceness *int     `yaml:"niceness",omitempty`
	Idler    Idler    `yaml:"idler"`
	Servers  []Server `yaml:"servers"`
}

type Server struct {
	Path      string `yaml:"path"`
	Number    int    `yaml:"number"`
	Ip        string `yaml:"ip"`
	Port      int    `yaml:"port"`
	Core      *int   `yaml:"core",omitempty`
	Niceness  *int   `yaml:"niceness",omitempty`
	Map       string `yaml:"map"`
	Params    string `yaml:"params"`
	ParamsAdd string `yaml:"params_add"`
	Idler     Idler  `yaml:"idler"`
}

type Idler struct {
	Enabled  *bool `yaml:"enabled",omitempty`
	Niceness *int  `yaml:"niceness",omitempty`
}

// Path current config file path
var Path = GetDefaultLocation()

// Value current config value
var Value = Root{}

var defaults = struct {
	Project       string
	Ip            string
	Niceness      int
	IdlerEnabled  bool
	IdlerNiceness int
}{
	Project:       "srcds-server",
	Ip:            "127.0.0.1",
	Niceness:      0,
	IdlerEnabled:  false,
	IdlerNiceness: -20,
}

// GetDefaultLocation returns default location and name for the config file (~/srcds.yaml)
func GetDefaultLocation() string {
	homedir, _ := os.LookupEnv("HOME")
	return filepath.Join(homedir, "/srcds.yaml")
}

// Read reads and validates config file from path
func Read() error {
	Reset()
	// check if file exists
	_, err := os.Stat(Path)
	if err != nil {
		return err
	}
	// read it
	f, err := os.ReadFile(Path)
	if err != nil {
		return err
	}
	// cook it
	cfg := Root{}
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return err
	}
	err = validate(Path, cfg)
	if err != nil {
		return err
	}
	// serve
	Value = *setDefaults(&cfg)
	return nil
	// ... delicious?
}

func ReadOrThrow() {
	err := Read()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}

// Reset resets Value to its unread state
func Reset() {
	Value = Root{}
	setDefaults(&Value)
}

func validate(path string, cfg Root) error {
	if len(cfg.Servers) == 0 {
		return errors.New(fmt.Sprintf("config file \"%s\" must include at least one server", path))
	}
	everyServerHasPath := true
	for i, server := range cfg.Servers {
		everyServerHasPath = everyServerHasPath && server.Path != ""
		if !everyServerHasPath && cfg.Path == "" {
			return errors.New(fmt.Sprintf("config file \"%s\" must define top-level `path`, since not every server has a `path` defined", path))
		}
		if server.Port == 0 {
			return errors.New(fmt.Sprintf("config file \"%s\" must have a `port` defined in `servers[%v]`", path, i))
		}
	}

	return nil
}

func setDefaults(cfg *Root) *Root {
	if cfg.Project == "" {
		cfg.Project = defaults.Project
	}
	if cfg.Ip == "" {
		cfg.Ip = defaults.Ip
	}
	if cfg.Idler.Enabled == nil {
		cfg.Idler.Enabled = &defaults.IdlerEnabled
	}
	if cfg.Idler.Niceness == nil {
		cfg.Idler.Niceness = &defaults.IdlerNiceness
	}
	if cfg.Niceness == nil {
		cfg.Niceness = &defaults.Niceness
	}
	// propagate top-level values to each server
	for i, server := range cfg.Servers {
		if server.Path == "" {
			server.Path = cfg.Path
		}
		if server.Ip == "" {
			server.Ip = cfg.Ip
		}
		if server.Map == "" {
			server.Map = cfg.Map
		}
		if server.Params == "" {
			server.Params = cfg.Params
		}
		if server.Number == 0 {
			server.Number = i + 1
		}
		if server.Niceness == nil {
			server.Niceness = cfg.Niceness
		}
		if server.Idler.Enabled == nil {
			server.Idler.Enabled = cfg.Idler.Enabled
		}
		if server.Idler.Niceness == nil {
			server.Idler.Niceness = cfg.Idler.Niceness
		}
		cfg.Servers[i] = server
	}

	return cfg
}
