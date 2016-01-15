// Package flagcfg populates flags from a TOML config file.
// Each flag is assumed to have an optional top-level value
// in the config file, having the same name. However, if a
// flag contains a dash or a period, those are converted to
// underscores.
//
// Flags that have aready been assigned are not overwritten.
//
// This package can be used together with github.com/facebookgo/flagenv
// to load flags from a config file, environment variable, or command-line.
package flagcfg

import (
	"errors"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	configLocator  Locator
	parsedFilename string
)

// ParseSet parses the configuration data in TOML format, placing found values into the flag set.
func ParseSet(tomlData []byte, set *flag.FlagSet) error {
	explicitlySet := make(map[string]bool)
	set.Visit(func(f *flag.Flag) {
		explicitlySet[f.Name] = true
	})

	var values map[string]interface{}
	if _, err := toml.Decode(string(tomlData), &values); err != nil {
		return errors.New(fmt.Sprintf("Unable to parse TOML: %s", err))
	}

	// log.Printf("%#v", values)

	var err error
	set.VisitAll(func(f *flag.Flag) {
		if err != nil {
			return
		}
		if _, ok := explicitlySet[f.Name]; !ok {
			name := strings.Replace(f.Name, ".", "_", -1)
			name = strings.Replace(name, "-", "_", -1)
			val, ok := values[name]
			var ferr error
			if ok {
				flagval := f.Value.(flag.Getter).Get()
				switch flagval.(type) {
				case string:
					ferr = f.Value.Set(val.(string))
				case time.Duration:
					// durations are just strings in TOML
					ferr = f.Value.Set(val.(string))
				case float64:
					ferr = f.Value.Set(strconv.FormatFloat(val.(float64), 'G', -1, 64))
				case int64, int:
					ferr = f.Value.Set(strconv.FormatInt(val.(int64), 10))
				case uint64, uint:
					ferr = f.Value.Set(strconv.FormatUint(uint64(val.(int64)), 10))
				case bool:
					ferr = f.Value.Set(strconv.FormatBool(val.(bool)))
				case flag.Value:
					ferr = f.Value.Set(val.(flag.Value).String())
				default:
					ferr = errors.New(fmt.Sprintf("Unable to map type of %#v", val))
				}
				if ferr != nil {
					err = fmt.Errorf("failed to set flag %q with value %q", f.Name, val)
				}
			}
		}
	})
	return err
}

// AddFile adds a location to search for a config file.
func AddFile(fileName string) {
	configLocator.AddFile(fileName)
}

// FindConfig returns the first config file that exists in the list.
func FindConfig() string {
	return configLocator.FindConfig()
}

// Parse will set each defined flag from the first configuration file
// found in the list of those added.
func Parse() {
	if configLocator.Len() == 0 {
		log.Fatalln("No configuration files specified")
	}

	parsedFilename = configLocator.FindConfig()
	if parsedFilename != "" {
		// log.Printf("Loading configuration from %s", parsedFilename)
		b, err := ioutil.ReadFile(parsedFilename)
		if err != nil {
			log.Fatalf("Unable to read %s: %s", parsedFilename, err)
		}
		if err := ParseSet(b, flag.CommandLine); err != nil {
			log.Fatalln(err)
		}
	}
}

// Filename returns the name of the config file that was parsed by Parse()
func Filename() string {
	return parsedFilename
}

// AddDefaults adds default config file locations, using the binary
// name as a base. It will look for the file specificed by
// {NAME}_CONFIG, as well as "{NAME}.config" in various locations.
func AddDefaults() {
	configLocator.AddDefaults()
}

// AddDefaultFiles adds default config file locations to search,
// using the given environment variable name and name of the
// config file to look for (excluding path).
func AddDefaultFiles(envName, cfgName string) {
	configLocator.AddDefaultFiles(envName, cfgName)
}
