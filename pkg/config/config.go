package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RepoForksConfig struct {
	Forks []RepoForkConfig `yaml:"forks"`
}

type RepoForkConfig struct {
	Name           string   `yaml:"name,omitempty"`
	OriginRepo     string   `yaml:"originRepo"`
	OriginBranch   string   `yaml:"originBranch,omitempty"`
	UpstreamRepo   string   `yaml:"upstreamRepo"`
	UpstreamBranch string   `yaml:"upstreamBranch,omitempty"`
	FullRewrite    bool     `yaml:"fullRewrite,omitempty"`
	Push           bool     `yaml:"push,omitempty"`
	ExcludePaths   []string `yaml:"excludePaths,omitempty"`
}

func LoadConfig(path string) (*RepoForksConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg RepoForksConfig
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// default values
	if cfg.Forks != nil {
		for i, fork := range cfg.Forks {
			if fork.OriginBranch == "" {
				cfg.Forks[i].OriginBranch = "main"
			}
			if fork.UpstreamBranch == "" {
				cfg.Forks[i].UpstreamBranch = "main"
			}
			if fork.ExcludePaths == nil {
				cfg.Forks[i].ExcludePaths = []string{
					".github/workflows/",
					".gitlab-ci.yml",
					".github/",
					".gitlab/",
				}
			}
		}
	}

	return &cfg, nil
}
