package main

import (
	"os"

	"github.com/jkrivas/k8s-ha-git-sync/internal/app/k8s-ha-git-sync"
	"github.com/jkrivas/k8s-ha-git-sync/internal/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	cfg := config.Config{}

	cmd := &cli.App{
		Name:            "k8s-ha-git-sync",
		Usage:           "Synchronize Kubernetes deployed Home Assistant configuration with Git repository",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "interval",
				Usage:       "Interval in seconds between synchronizations",
				Value:       60,
				EnvVars:     []string{"INTERVAL"},
				Destination: &cfg.Interval,
			},
			&cli.StringFlag{
				Name:        "ha-config-path",
				Category:    "Home Assistant",
				Usage:       "Path to the Home Assistant configuration directory",
				Value:       "/homeassistant",
				EnvVars:     []string{"CONFIG_PATH"},
				Destination: &cfg.HomeAssistant.ConfigPath,
			},
			&cli.StringFlag{
				Name:        "ha-url",
				Category:    "Home Assistant",
				Usage:       "URL of the Home Assistant instance",
				Value:       "http://homeassistant:8123",
				EnvVars:     []string{"HA_URL"},
				Destination: &cfg.HomeAssistant.Url,
			},
			&cli.StringFlag{
				Name:        "ha-token",
				Category:    "Home Assistant",
				Usage:       "Long-Lived Access Token for the Home Assistant instance",
				EnvVars:     []string{"HA_TOKEN"},
				Required:    true,
				Destination: &cfg.HomeAssistant.Token,
			},
			&cli.StringFlag{
				Name:        "git-ssh-key-path",
				Category:    "Git",
				Usage:       "Path to the SSH key for Git authentication",
				EnvVars:     []string{"GIT_SSH_KEY_PATH"},
				Destination: &cfg.Git.SshKeyPath,
			},
			&cli.StringFlag{
				Name:        "git-token",
				Category:    "Git",
				Usage:       "Token for Git authentication",
				EnvVars:     []string{"GIT_TOKEN"},
				Destination: &cfg.Git.Token,
			},
			&cli.StringFlag{
				Name:        "kube-namespace",
				Category:    "Kubernetes",
				Usage:       "Name of the Home Assistant deployment namespace in Kubernetes",
				Value:       "default",
				EnvVars:     []string{"KUBE_NAMESPACE"},
				Destination: &cfg.Kubernetes.Namespace,
			},
			&cli.StringFlag{
				Name:        "kube-deployment",
				Category:    "Kubernetes",
				Usage:       "Name of the Home Assistant deployment in Kubernetes",
				Value:       "homeassistant",
				EnvVars:     []string{"KUBE_DEPLOYMENT"},
				Destination: &cfg.Kubernetes.Deployment,
			},
			&cli.BoolFlag{
				Name:        "metrics",
				Category:    "Metrics",
				Usage:       "Enable Prometheus metrics",
				Value:       false,
				EnvVars:     []string{"METRICS"},
				Destination: &cfg.Metrics.Enabled,
			},
			&cli.IntFlag{
				Name:        "metrics-port",
				Category:    "Metrics",
				Usage:       "Port for Prometheus metrics",
				Value:       8080,
				EnvVars:     []string{"METRICS_PORT"},
				Destination: &cfg.Metrics.Port,
			},
		},
		Before: func(c *cli.Context) error {
			if c.String("git-ssh-key-path") == "" && c.String("git-token") == "" {
				log.Fatal("one of --git-ssh-key-path or --git-token must be provided")
			}

			if c.String("git-ssh-key-path") != "" && c.String("git-token") != "" {
				log.Fatal("only one of --git-ssh-key-path or --git-token must be provided")
			}
			return nil
		},
		Action: func(ctx *cli.Context) error {
			sync.Sync(&cfg)
			return nil
		},
	}

	if err := cmd.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
