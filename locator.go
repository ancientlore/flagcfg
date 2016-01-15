package flagcfg

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Locator can be used to find a given configuration file from a list you set up.
type Locator struct {
	configFiles []string
}

// Len returns the number of config files locations that have been added.
func (c *Locator) Len() int {
	if c.configFiles == nil {
		return 0
	}
	return len(c.configFiles)
}

// Reset resets the locator to be used again.
func (c *Locator) Reset() {
	c.configFiles = nil
}

// AddFile adds a location to search for a config file.
func (c *Locator) AddFile(fileName string) {
	if c.configFiles == nil {
		c.configFiles = make([]string, 0, 8)
	}

	c.configFiles = append(c.configFiles, fileName)
}

// FindConfig returns the first config file that exists in the list.
func (c *Locator) FindConfig() string {
	if c.configFiles != nil {
		for _, filename := range c.configFiles {
			_, err := os.Stat(filename)
			if err == nil {
				return filename
			}
		}
	}
	return ""
}

func getExePath() (string, string, string) {
	var exePath, err = exec.LookPath(os.Args[0])
	if err != nil {
		log.Print("Warning: ", err)
		exePath = os.Args[0]
	}
	s, err := filepath.Abs(exePath)
	if err != nil {
		log.Print("Warning: ", err)
	} else {
		exePath = s
	}
	p, n := filepath.Split(exePath)
	e := filepath.Ext(n)
	n = strings.TrimSuffix(n, e)
	return p, n, e
}

// AddDefaults adds default config file locations, using the binary
// name as a base. It will look for the file specificed by
// {NAME}_CONFIG, as well as "{NAME}.config" in various locations.
func (c *Locator) AddDefaults() {
	_, n, _ := getExePath()
	c.AddDefaultFiles(strings.ToUpper(n)+"_CONFIG", n+".config")
}

// AddDefaultFiles adds default config file locations to search,
// using the given environment variable name and name of the
// config file to look for (excluding path).
func (c *Locator) AddDefaultFiles(envName, cfgName string) {
	// look at environment variable first
	if envName != "" {
		cfgFile := os.Getenv(envName)
		if cfgFile != "" {
			c.AddFile(cfgFile)
		}
	}

	// break down EXE name
	exePath, exeName, _ := getExePath()

	// current folder
	wd, err := os.Getwd()
	if err == nil {
		c.AddFile(filepath.Join(wd, cfgName))
	}

	// Look in home folder
	home := os.Getenv("HOME")
	if home != "" {
		c.AddFile(filepath.Join(home, ".config", exeName, cfgName))
	}

	// some etc folder relative to binary, assuming binary isn't directly in /bin
	if strings.Contains(exePath, filepath.FromSlash("/bin/")) && exePath != filepath.FromSlash("/bin") {
		c.AddFile(filepath.Join(strings.Replace(exePath, filepath.FromSlash("/bin/"), filepath.FromSlash("/etc/"), 1), cfgName))
	}

	// etc locations
	c.AddFile(filepath.Join(filepath.FromSlash("/etc"), exeName, cfgName))
	c.AddFile(filepath.Join(filepath.FromSlash("/etc"), cfgName))

	// Same folder as binary
	c.AddFile(filepath.Join(exePath, cfgName))

	// log.Print(c.configFiles)
}
