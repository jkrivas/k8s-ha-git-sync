package git

import (
	"fmt"
	"strings"

	"github.com/gogs/git-module"
)

func Open(p string) (*git.Repository, error) {
	return git.Open(p)
}

func PullRepo(r *git.Repository) (bool, error) {
	crev, err := r.RevParse("HEAD")
	if err != nil {
		return false, err
	}

	err = r.Pull()
	if err != nil {
		return false, err
	}

	nrev, err := r.RevParse("HEAD")
	if err != nil {
		return false, err
	}

	if crev == nrev {
		return false, nil
	}

	return true, nil
}

func getRemote(r *git.Repository) (string, error) {
	rem, err := r.Remotes()
	if err != nil {
		return "", err
	}

	if len(rem) == 0 {
		return "", fmt.Errorf("no remote found")
	}

	return rem[0], nil
}

func getRemoteURL(r *git.Repository) (string, error) {
	rem, err := getRemote(r)
	if err != nil {
		return "", err
	}

	url, err := r.RemoteGetURL(rem)
	if err != nil {
		return "", err
	}

	if len(url) == 0 {
		return "", fmt.Errorf("no remote URL found")
	}

	return url[0], nil
}

func getRemoteProtocol(r *git.Repository) (string, error) {
	u, err := getRemoteURL(r)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(u, "git@") {
		return "ssh", nil
	} else if strings.HasPrefix(u, "https://") {
		return "https", nil
	}

	return "", fmt.Errorf("git origin protocol not supported: %s", u)
}
