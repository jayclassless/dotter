package step

import "os"

// IsSymLink indicates whether or not the speicifed FileInfo describes a Symlink
func IsSymLink(fileInfo os.FileInfo) bool {
	return fileInfo.Mode()&os.ModeSymlink != 0
}
