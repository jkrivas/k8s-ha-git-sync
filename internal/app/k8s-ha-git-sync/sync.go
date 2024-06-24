package sync

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	g "github.com/gogs/git-module"
	k "k8s.io/client-go/kubernetes"

	"github.com/jkrivas/k8s-ha-git-sync/internal/pkg/config"
	"github.com/jkrivas/k8s-ha-git-sync/internal/pkg/git"
	"github.com/jkrivas/k8s-ha-git-sync/internal/pkg/homeassistant"
	"github.com/jkrivas/k8s-ha-git-sync/internal/pkg/kubernetes"
)

var (
	configStatus = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ha_git_sync_config_status",
		Help: "Status of the Home Assistant configuration",
	})
)

func sync(cfg *config.Config, kcs *k.Clientset, gr *g.Repository) {
	for {
		log.Info("syncing...")
		change, err := git.PullRepo(gr)
		if err != nil {
			log.Error(err)
		}

		err = homeassistant.CheckConfig(cfg.HomeAssistant.Url, cfg.HomeAssistant.Token)
		if err != nil {
			log.Error(err)
			configStatus.Set(0)
		} else {
			configStatus.Set(1)
		}

		if change && err == nil {
			log.Info("validated new config")
			err = kubernetes.RestartDeployment(kcs, cfg.Kubernetes.Namespace, cfg.Kubernetes.Deployment)
			if err != nil {
				log.Error(err)
			} else {
				log.Infof("deployment %v:%v restarted", cfg.Kubernetes.Namespace, cfg.Kubernetes.Deployment)
			}
		}

		time.Sleep(time.Duration(cfg.Interval) * time.Second)
	}
}

func Sync(cfg *config.Config) {
	log.Info("starting app...")

	kc, err := kubernetes.GenerateClient(cfg.Kubernetes.KubeconfigPath)
	if err != nil {
		log.Fatal(err)
	}

	gr, err := git.Open(cfg.HomeAssistant.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	err = git.Authenticate(gr, cfg.Git.SshKeyPath, cfg.Git.Token)
	if err != nil {
		log.Fatal(err)
	}

	go sync(cfg, kc, gr)

	if cfg.Metrics.Enabled {
		log.Info("starting metrics server...")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Metrics.Port), nil); err != nil {
			log.Fatal(err)
		}
	}

	// Block forever
	select {}
}
