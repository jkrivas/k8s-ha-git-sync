package git

import (
	"fmt"
	"github.com/gogs/git-module"
	"github.com/jkrivas/k8s-ha-git-sync/internal/pkg/utils"
)

func authenticateWithToken(token string) error {
	if token == "" {
		return fmt.Errorf("no token provided")
	}

	cmd := git.NewCommand("config", "--global", "credential.helper", fmt.Sprintf("!f() { echo username=token; echo \"password=%s\"; }; f", token))
	_, err := cmd.Run()
	return err
}

func authenticateWithSSHKey(sshKeyPath string) error {
	if sshKeyPath == "" {
		return fmt.Errorf("no SSH key path provided")
	}

	var p = "/tmp/ssh_key_git"
	if err := utils.CopyFile(sshKeyPath, p, 0600); err != nil {
		return err
	}

	cmd := git.NewCommand("config", "--global", "core.sshCommand", fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes -o StrictHostKeychecking=no", p))
	if _, err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func Authenticate(gr *git.Repository, sshKeyPath string, token string) error {
	rp, err := getRemoteProtocol(gr)
	if err != nil {
		return err
	}

	switch rp {
	case "https":
		return authenticateWithToken(token)
	case "ssh":
		return authenticateWithSSHKey(sshKeyPath)
	}

	return fmt.Errorf("unsupported git protocol: %s", rp)
}
