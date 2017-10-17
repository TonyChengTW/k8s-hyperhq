// +build linux

package overlay

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"syscall"

	"github.com/hyperhq/hyperd/utils"
)

func MountContainerToSharedDir(containerId, rootDir, sharedDir, mountLabel string, readonly bool) (string, error) {
	var (
		params     string
		mountPoint = path.Join(sharedDir, containerId, "rootfs")
		upperDir   = path.Join(rootDir, containerId, "upper")
		workDir    = path.Join(rootDir, containerId, "work")
	)

	if _, err := os.Stat(mountPoint); err != nil {
		if err = os.MkdirAll(mountPoint, 0755); err != nil {
			return "", err
		}
	}
	lowerId, err := ioutil.ReadFile(path.Join(rootDir, containerId) + "/lower-id")
	if err != nil {
		return "", err
	}
	lowerDir := path.Join(rootDir, string(lowerId), "root")

	if readonly {
		// "upperdir=" and "workdir=" may be omitted. In that case the overlay will be read-only.
		params = fmt.Sprintf("lowerdir=%s:%s", lowerDir, upperDir)
	} else {
		params = fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir)
	}
	if err := syscall.Mount("overlay", mountPoint, "overlay", 0, utils.FormatMountLabel(params, mountLabel)); err != nil {
		return "", fmt.Errorf("error creating overlay mount to %s: %v", mountPoint, err)
	}
	return mountPoint, nil
}
