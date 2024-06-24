# k8s-ha-git-sync

Toll that allows to sync Kubernetes deployed Home Assistant configuration with Git.

Inspired by [Home Assistant Git Pull addon](https://github.com/home-assistant/addons/tree/master/git_pull).

## Introduction

The tool works by periodically executing a `git pull` command (default interval: 60s). It then checks the validity of the configuration by calling the Home Assistant API. If the configuration is valid, it triggers a deployment restart using the Kubernetes API.

> [!WARNING]
> The developer of this tool takes no responsibility for any unexpected changes or deletions to your Home Assistant configuration. It is your responsibility to ensure that you have a backup of your configuration before using this tool.

## Pre-setup

### Home Assistant

In order to validate configuration in Home Assistant, `api` integration needs to be enabled: [documentation](https://www.home-assistant.io/integrations/api/).

Then in users profile a long-lived access token for this tool needs to be generated.

### Git

Before deploying this tool you must correctly setup the Home Assistant configuration directory as a Git repository with a remote repository.
Currently the tool does not support initializing a git repository, or cloning one.

Official Home Assistant container image has git tool pre-installed, and can be used by getting a shell in the pod.

```sh
kubectl -n <namespace> exec --stdin --tty <pod> -- /bin/bash
cd /config 
git init -b <branch>
git remote add origin <repo>
git branch --set-upstream-to=origin/<branch> <branch>
```

#### .gitignore

Since Home Assistants keeps other various service files like logs, databases, backups etc... , it is recommended to ignore everything in `.gitignore` file and then only allow synced files.

Example `.gitignore`:

```text
# Ignore everything
/*

# Synced files
!.gitignore
!configuration.yaml
!configuration/
```

## Deployment

### Options

| Option             | Environment variable | Description                                                   | Default                  | Required |
| ------------------ | -------------------- | ------------------------------------------------------------- | ------------------------ | -------- |
| --interval         | INTERVAL             | Interval in seconds between synchronizations                  | 60                       | Yes      |
| --ha-config-path   | CONFIG_PATH          | Path to the Home Assistant configuration directory            | /homeassistant           | Yes      |
| --ha-url           | HA_URL               | URL of the Home Assistant instance                            | http://homeassitant:8123 | Yes      |
| --ha-token         | HA_TOKEN             | Long-Lived Access Token for the Home Assistant instance       |                          | Yes      |
| --git-ssh-key-path | GIT_SSH_KEY_PATH     | Path to the SSH key for Git authentication                    |                          | No       |
| --git-token        | GIT_TOKEN            | Token for Git HTTPS authentication                            |                          | No       |
| --kube-namespace   | KUBE_NAMEPSACE       | Name of the Home Assistant deployment namespace in Kubernetes | homeassistant            | Yes      |
| --kube-deployment  | KUBE_DEPLOYMENT      | Name of the Home Assistant deployment in Kubernetes           | homeassistant            | Yes      |
| --metrics          | METRICS              | Enable Prometheus metrics                                     | false                    | No       |
| --metrics-port     | METRICS_PORT         | Port for Prometheus metrics service                           | 8080                     | No       | ~~~~ |

### Kubernetes service account

Since this tool uses native Kubernetes API, it uses in-cluster authentication with permissions from the service account of the pod.
In order for it to restart deployments, a role and a role binding needs to be created.

Example:

```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: homeassistant-sync
  namespace: homeassistant
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: edit-deployments
  namespace: homeassistant
rules:
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get", "update"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: edit-deployments-homeassistant-sync
  namespace: homeassistant
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: edit-deployments
subjects:
  - kind: ServiceAccount
    name: homeassistant-sync
    namespace: homeassistant
```

## Observability

The tool has capability to expose Prometheus metrics. It can be enabled by setting env variable `METRICS` to `true`.

Exported metrics:

| Metric                    | Type  | Description                                             |
| ------------------------- | ----- | ------------------------------------------------------- |
| ha_git_sync_config_status | Gauge | Shows if pulled configuration is valid. Returns 1 or 0. |
