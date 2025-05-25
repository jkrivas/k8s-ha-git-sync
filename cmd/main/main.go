package main

import (
	"context"
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

	cmd := &cli.Command{
		Name:            "k8s-ha-git-sync",
		Usage:           "Synchronize Kubernetes deployed Home Assistant configuration with Git repository",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "interval",
				Usage:       "Interval in seconds between synchronizations",
				Value:       60,
				Sources:     cli.EnvVars("INTERVAL"),
				Destination: &cfg.Interval,
			},
			&cli.StringFlag{
				Name:        "ha-config-path",
				Category:    "Home Assistant",
				Usage:       "Path to the Home Assistant configuration directory",
				Value:       "/homeassistant",
				Sources:     cli.EnvVars("CONFIG_PATH"),
				Destination: &cfg.HomeAssistant.ConfigPath,
			},
			&cli.StringFlag{
				Name:        "ha-url",
				Category:    "Home Assistant",
				Usage:       "URL of the Home Assistant instance",
				Value:       "http://homeassistant:8123",
				Sources:     cli.EnvVars("HA_URL"),
				Destination: &cfg.HomeAssistant.Url,
			},
			&cli.StringFlag{
				Name:        "ha-token",
				Category:    "Home Assistant",
				Usage:       "Long-Lived Access Token for the Home Assistant instance",
				Sources:     cli.EnvVars("HA_TOKEN"),
				Required:    true,
				Destination: &cfg.HomeAssistant.Token,
			},
			&cli.StringFlag{
				Name:        "git-ssh-key-path",
				Category:    "Git",
				Usage:       "Path to the SSH key for Git authentication",
				Sources:     cli.EnvVars("GIT_SSH_KEY_PATH"),
				Destination: &cfg.Git.SshKeyPath,
			},
			&cli.StringFlag{
				Name:        "git-token",
				Category:    "Git",
				Usage:       "Token for Git authentication",
				Sources:     cli.EnvVars("GIT_TOKEN"),
				Destination: &cfg.Git.Token,
			},
			&cli.StringFlag{
				Name:        "kube-namespace",
				Category:    "Kubernetes",
				Usage:       "Name of the Home Assistant deployment namespace in Kubernetes",
				Value:       "default",
				Sources:     cli.EnvVars("KUBE_NAMESPACE"),
				Destination: &cfg.Kubernetes.Namespace,
			},
			&cli.StringFlag{
				Name:        "kube-deployment",
				Category:    "Kubernetes",
				Usage:       "Name of the Home Assistant deployment in Kubernetes",
				Value:       "homeassistant",
				Sources:     cli.EnvVars("KUBE_DEPLOYMENT"),
				Destination: &cfg.Kubernetes.Deployment,
			},
			&cli.BoolFlag{
				Name:        "metrics",
				Category:    "Metrics",
				Usage:       "Enable Prometheus metrics",
				Value:       false,
				Sources:     cli.EnvVars("METRICS"),
				Destination: &cfg.Metrics.Enabled,
			},
			&cli.IntFlag{
				Name:        "metrics-port",
				Category:    "Metrics",
				Usage:       "Port for Prometheus metrics",
				Value:       8080,
				Sources:     cli.EnvVars("METRICS_PORT"),
				Destination: &cfg.Metrics.Port,
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			if cmd.String("git-ssh-key-path") == "" && cmd.String("git-token") == "" {
				log.Fatal("one of --git-ssh-key-path or --git-token must be provided")
			}

			if cmd.String("git-ssh-key-path") != "" && cmd.String("git-token") != "" {
				log.Fatal("only one of --git-ssh-key-path or --git-token must be provided")
			}
			return nil, nil
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			sync.Sync(&cfg)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
