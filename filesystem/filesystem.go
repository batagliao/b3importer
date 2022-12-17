package filesystem

import "os"

func FileExists(file string) bool {
	_, err := os.Stat(file)

	return !os.IsNotExist(err)
}

func DeleteFile(file string) error {
	return os.Remove(file)
}
