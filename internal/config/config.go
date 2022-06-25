package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

type Config struct {
	HttpPort         int           `yaml:"http_port"`
	PingInterval     time.Duration `yaml:"ping_interval"`
	HistogramBuckets []float64     `yaml:"histogram_buckets"`
	LogsEnv          string        `yaml:"logs_env"`

	HostId      string                 `yaml:"host_id"`
	Probes      map[string]*Probe      `yaml:"probes"`
	Deployments map[string]*Deployment `yaml:"deployments"`
}

type Probe struct {
	Host    string   `yaml:"host"`
	Targets []string `yaml:"targets"`
}

type Deployment struct {
	Hosts map[string]*Host `yaml:"hosts"`
}

type Host struct {
	Address string `yaml:"address"`
}

func Load(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	//zlog.Info().Str(logfields.Module, "config").Str("config", string(file)).Send()

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
