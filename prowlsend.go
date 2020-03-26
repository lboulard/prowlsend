package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/rem7/goprowl"
)

var Options struct {
	Priority    int
	Event       string
	Url         string
	Application string
	HostPrefix  bool
	ConfigFile  string
}

const CONFIG_PATH = "prowl/prowl.toml"

type Config struct {
	ApiKey string
}

func ReadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func ConfigFileLocations(configpath string) []string {
	locations := []string{}

	if runtime.GOOS == "windows" {
		configdirs := []string{}
		appdata := os.Getenv("APPDATA")
		if len(appdata) > 0 {
			configdirs = append(configdirs, appdata)
		}
		localappdata := os.Getenv("LOCALAPPDATA")
		if len(localappdata) > 0 {
			configdirs = append(configdirs, localappdata)
		}
		for _, configdir := range configdirs {
			locations = append(locations, filepath.Join(configdir, configpath))
		}
	} else {
		for _, path := range strings.Split(os.Getenv("XDG_CONFIG_DIRS"), ":") {
			if len(path) > 0 {
				locations = append(locations, filepath.Join(path, configpath))
			}
		}
		if len(locations) == 0 {
			user, _ := user.Current()
			if user != nil && user.HomeDir != "" {
				locations = append(locations,
					filepath.Join(user.HomeDir, ".config", configpath))
			}
		}
	}
	return locations
}

func FindConfigFile(configpath string) (string, error) {
	var locations []string
	if len(Options.ConfigFile) > 0 {
		locations = []string{Options.ConfigFile}
	} else {
		locations = ConfigFileLocations(configpath)
	}
	for _, path := range locations {
		if fi, err := os.Stat(path); err == nil {
			if fi.Mode().IsRegular() {
				perm := fi.Mode().Perm()
				if (perm & 0027) != 0 {
					fmt.Fprintf(os.Stderr,
						"** WARNING: file permission for \"%s\" are too large\n",
						path)
				}
				return path, nil
			}
		}
	}

	var err error
	if len(locations) > 0 {
		err = fmt.Errorf("%#v not found", filepath.Join(locations[0], configpath))
	} else {
		err = fmt.Errorf("Configuration path must be given with `-c` option")
	}
	return "", err
}

func main() {
	flag.Parse()

	configpath, err := FindConfigFile(CONFIG_PATH)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	config, err := ReadConfig(configpath)
	if config == nil {
		fmt.Fprintf(os.Stderr, "Error:%s: %v\n", configpath, err)
		os.Exit(1)
	}
	if len(config.ApiKey) == 0 {
		fmt.Fprintf(os.Stderr, "Error: API key not found in configuration file\n")
		os.Exit(1)
	}

	var p goprowl.Goprowl
	if err := p.RegisterKey(config.ApiKey); err != nil {
		fmt.Fprintf(os.Stderr, "Error registering key:  %v\n", err)
		os.Exit(1)
	}

	app := Options.Application
	if Options.HostPrefix {
		if hostname, ok := os.Hostname(); ok == nil {
			app = app + " on " + hostname
		}
	}

	n := &goprowl.Notification{
		Application: app,
		Description: "",
		Event:       Options.Event,
		Priority:    strconv.Itoa(Options.Priority),
		Url:         Options.Url,
	}

	switch len(flag.Args()) {
	case 0:
		fmt.Fprintln(os.Stderr, "Error: need at least a message as argument")
		os.Exit(1)
	default:
		n.Description = strings.Join(flag.Args(), "\n")
	}

	err = p.Push(n)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: failed to send Prowl notification")
		fmt.Fprintln(os.Stderr, "# %r", err)
		os.Exit(1)
	}
}

func usage() {
	locations := ConfigFileLocations(CONFIG_PATH)
	fmt.Fprintf(os.Stderr, "Usage: %s [options] \"message\" [...]\n",
		filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "%s", `
Messages are concatenated with a carriage return character between them.
`)
	if len(locations) > 1 {
		fmt.Fprintf(os.Stderr, "%s%s\n", `
Default configuration files can be found at:
  - `, strings.Join(locations, "\n  - "))
	} else if len(locations) == 1 {
		fmt.Fprintf(os.Stderr, "%s%s\n", `
Default configuration file can be found at:
  - `, locations[0])
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", `
Cannot find a default configuration file location.
You must provide a configuration file with '-c' option.`)
	}
	fmt.Fprintf(os.Stderr, `
Options:
`)
	flag.PrintDefaults()
}

func init() {
	flag.IntVar(&Options.Priority, "p", 1, "notification `priority`")
	flag.StringVar(&Options.Event, "e", "", "notification optional `event`")
	flag.StringVar(&Options.Url, "u", "", "notification optional `URL`")
	flag.StringVar(&Options.Application, "a", "Prowlsend", "application `name`")
	flag.BoolVar(&Options.HostPrefix, "o", true,
		"Use \"-o=false\" to not append hostname in application name")
	flag.StringVar(&Options.ConfigFile, "c", "", "`path` to specific configuration file")
	flag.Usage = usage
}
