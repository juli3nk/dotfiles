package helpers

import (
	"os"

	"github.com/juli3nk/dotfiles/internal/config"
	"github.com/juli3nk/dotfiles/internal/dotfiles"
	"github.com/juli3nk/dotfiles/internal/registry"
)

type Composed struct {
	Reg        *registry.Registry
	Dots       []*dotfiles.Dotfiles
	DotConfigs []*config.Config
	Config     *config.Config
}

func OpenRegistry() (*registry.Registry, error) {
	return registry.New(os.Getenv("HOME"))
}

// OpenComposed loads the registry and builds the composed set of repositories
// (default/base first, then layer repos in registration order) together with
// their merged configuration. DotConfigs[i] is the individual config of
// Dots[i]; Config is the merged config across all composed repositories.
func OpenComposed() (*Composed, error) {
	reg, err := OpenRegistry()
	if err != nil {
		return nil, err
	}

	entries, err := reg.Composed()
	if err != nil {
		return nil, err
	}

	result := &Composed{Reg: reg}

	for _, e := range entries {
		d, err := dotfiles.NewWithPath(e.Name, reg.CloneDir(e.Name))
		if err != nil {
			return nil, err
		}

		cfg, err := d.LoadConfig()
		if err != nil {
			return nil, err
		}

		result.Dots = append(result.Dots, d)
		result.DotConfigs = append(result.DotConfigs, cfg)

		if result.Config == nil {
			result.Config = cfg
		} else {
			result.Config.Merge(cfg)
		}
	}

	return result, nil
}
