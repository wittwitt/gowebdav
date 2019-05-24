package server

import (
	"os"
	"path/filepath"
)

// GetInstancePath get instance path
func GetInstancePath() (string, error) {
	argPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Dir(argPath), nil

}

//PathExist check  path exists
func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
