package pinger

import (
	"context"
	"github.com/auntan/network_exporter/internal/metrics"
	"github.com/go-ping/ping"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"time"
)

type pingConfig struct {
	MyHost        string
	TargetHost    string
	TargetAddress string
	Interval      time.Duration

	MetricsConfig []*metricsConfig
}

type metricsConfig struct {
	Probe            string
	HostDeployment   string
	Host             string
	TargetDeployment string
	TargetHost       string
	TargetAddress    string
}

type singlePinger struct {
	conf *pingConfig

	pinger *ping.Pinger
}

func newSinglePinger(cfg *pingConfig) *singlePinger {
	return &singlePinger{
		conf: cfg,
	}
}

func (p *singlePinger) Run(ctx context.Context) error {
	pinger, err := ping.NewPinger(p.conf.TargetAddress)
	if err != nil {
		return err
	}
	p.pinger = pinger
	pinger.Interval = p.conf.Interval
	pinger.SetPrivileged(true)

	go func() {
		<-ctx.Done()
		pinger.Stop()
	}()

	pinger.OnRecv = p.onRecv

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(pinger.Run)

	eg.Go(func() error {
		var prevSent, prevRecv int
		for range time.Tick(time.Second) {
			if ctx.Err() != nil {
				return nil
			}

			stat := pinger.Statistics()

			sent := stat.PacketsSent - prevSent
			prevSent = stat.PacketsSent

			metrics.SentHost.WithLabelValues(
				p.conf.MyHost,
				p.conf.TargetHost,
				p.conf.TargetAddress,
			).Add(float64(sent))

			recv := stat.PacketsRecv - prevRecv
			prevRecv = stat.PacketsRecv

			metrics.RecvHost.WithLabelValues(
				p.conf.MyHost,
				p.conf.TargetHost,
				p.conf.TargetAddress,
			).Add(float64(recv))

			for _, m := range p.conf.MetricsConfig {
				metrics.SentProbe.WithLabelValues(
					m.Probe,
					m.HostDeployment,
					m.Host,
					m.TargetDeployment,
					m.TargetHost,
					m.TargetAddress,
				).Add(float64(sent))

				metrics.RecvProbe.WithLabelValues(
					m.Probe,
					m.HostDeployment,
					m.Host,
					m.TargetDeployment,
					m.TargetHost,
					m.TargetAddress,
				).Add(float64(recv))
			}
		}

		return nil
	})

	err = eg.Wait()

	return err
}

func (p *singlePinger) onRecv(pkt *ping.Packet) {
	zap.S().Infof("%d bytes from %s: icmp_seq=%d time=%v", pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)

	metrics.RTTHost.WithLabelValues(
		p.conf.MyHost,
		p.conf.TargetHost,
		p.conf.TargetAddress,
	).Observe(pkt.Rtt.Seconds())

	for _, m := range p.conf.MetricsConfig {
		metrics.RTTProbe.WithLabelValues(
			m.Probe,
			m.HostDeployment,
			m.Host,
			m.TargetDeployment,
			m.TargetHost,
			m.TargetAddress,
		).Observe(pkt.Rtt.Seconds())
	}
}
