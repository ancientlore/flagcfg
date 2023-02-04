flagcfg
=======

[![Go Reference](https://pkg.go.dev/badge/github.com/ancientlore/flagcfg.svg)](https://pkg.go.dev/github.com/ancientlore/flagcfg)

The flagcfg package populates flags from a TOML config file.
Each flag is assumed to have an optional top-level value
in the config file, having the same name. However, if a
flag contains a dash or a period, those are converted to
underscores.

Flags that have aready been assigned are not overwritten.

This package can be used together with github.com/facebookgo/flagenv
to load flags from a config file, environment variable, or command-line.

Example:

	// Parse flags from command-line
	flag.Parse()

	// Parser flags from config
	flagcfg.AddDefaults()
	// or use flagcfg.AddDefaultFiles("MYAPP_CONFIG", "myapp.config")
	flagcfg.Parse()

	// Parse flags from environment (using github.com/facebookgo/flagenv)
	flagenv.Prefix = "MYPREFIX_"
	flagenv.Parse()
