package helpers

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"minecraft_prelauncher/settings"
	"os"
	"path"
	"path/filepath"
	"strings"

	minio "github.com/minio/minio-go"
)

func IsMaintenance(client *minio.Client, bucket string) (bool, error) {
	doneCh := make(chan struct{})

	defer close(doneCh)

	for object := range client.ListObjects(bucket, "maintenance", true, doneCh) {
		if object.Err != nil {
			return false, object.Err
		}
		return true, nil
	}
	return false, nil
}

func ReadModPaths(root string) ([]string, error) {
	var files []string

	fileInfos, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, info := range fileInfos {
		if info.IsDir() || strings.HasSuffix(info.Name(), ".DS_Store") {
			continue
		}
		files = append(files, filepath.Join(root, info.Name()))
	}
	return files, nil
}

func CalculateChecksum(filepath string) string {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	checksum := sha256.Sum256(file)
	return fmt.Sprintf("%x", checksum)
}

func FileExists(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	}
	return false
}

func GenerateConfig() {
	if !FileExists("config.json") {
		settings.CreateConfigTemplate()
		fmt.Println("Config files has been created. Please fill them with required info and run the program again.")
		os.Exit(0)
	}
}

func ModPathsCorrect(modPaths []string) bool {
	fmt.Printf("Found %d local mods.\n", len(modPaths))
	for _, el := range modPaths {
		if strings.HasSuffix(el, ".jar") {
			continue
		}
		return false
	}
	return true
}

func IsFileDeprecated(manifest settings.Manifest, filepath string, checksum string, excludes []string) bool {
	for _, mod := range manifest.Records {
		if mod.Checksum == checksum {
			return false
		}
	}
	for _, excl := range excludes {
		if strings.Contains(path.Base(filepath), excl) {
			return false
		}
	}
	return true
}

func IsNewRemoteFile(record settings.Record, localChecksums []string) bool {
	if record.Name == "manifest.json" {
		return false
	}
	for _, el := range localChecksums {
		if record.Checksum == el {
			return false
		}
	}
	return true
}
