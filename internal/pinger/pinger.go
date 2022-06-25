package pinger

import (
	"context"
	"fmt"
	"github.com/auntan/network_exporter/internal/config"
	"github.com/auntan/network_exporter/internal/graceful"
)

type Pinger struct {
	conf *config.Config
}

func New(conf *config.Config) *Pinger {
	return &Pinger{
		conf: conf,
	}
}

func (p *Pinger) Run(ctx context.Context) error {
	pingConfigs := map[string]*pingConfig{}
	for probeId, probe := range p.conf.Probes {
		hostDeployment := p.conf.Deployments[probe.Host]
		if hostDeployment == nil {
			return fmt.Errorf("probe: %s; host deployment %s not found", probeId, probe.Host)
		}

		found := false
		for hostId := range hostDeployment.Hosts {
			if hostId == p.conf.HostId {
				found = true
				break
			}
		}
		if !found {
			// check next probe
			continue
		}

		for _, targetDeploymentId := range probe.Targets {
			targetDeployment := p.conf.Deployments[targetDeploymentId]
			if targetDeployment == nil {
				return fmt.Errorf("probe: %s; target deployment %s not found", probeId, targetDeploymentId)
			}

			for targetHostId, targetHost := range targetDeployment.Hosts {
				key := targetHostId
				if pingConfigs[key] == nil {
					pingConfigs[key] = &pingConfig{
						MyHost:        p.conf.HostId,
						TargetHost:    targetHostId,
						TargetAddress: targetHost.Address,
						Interval:      p.conf.PingInterval,
						MetricsConfig: []*metricsConfig{},
					}
				}
				pingConfigs[key].MetricsConfig = append(pingConfigs[key].MetricsConfig, &metricsConfig{
					Probe:            probeId,
					HostDeployment:   probe.Host,
					Host:             p.conf.HostId,
					TargetDeployment: targetDeploymentId,
					TargetHost:       targetHostId,
					TargetAddress:    targetHost.Address,
				})
			}
		}
	}

	if len(pingConfigs) == 0 {
		return fmt.Errorf("empty targets")
	}

	pings := map[string]func(context.Context) error{}
	for k, target := range pingConfigs {
		p := newSinglePinger(target)
		pings[fmt.Sprintf("ping_%s", k)] = p.Run
	}

	return graceful.Run(ctx, pings)
}
