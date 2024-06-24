package utils

import (
	"io"
	"os"
)

func CopyFile(src, dst string, perm uint32) error {
	srcf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcf.Close()

	dstf, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstf.Close()

	_, err = io.Copy(dstf, srcf)
	if err != nil {
		return err
	}

	// Set the permissions of the destination file
	err = dstf.Chmod(os.FileMode(perm))
	if err != nil {
		return err
	}

	return nil
}
