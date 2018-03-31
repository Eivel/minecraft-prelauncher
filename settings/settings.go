package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const manifestFilename = "manifest.json"
const configFilename = "config.json"

type Manifest struct {
	Records []Record `json:"records"`
}

type Record struct {
	Name     string `json:"name"`
	Checksum string `json:"checksum"`
}

type Config struct {
	Excluded     []string `json:"excluded"`
	ModsPath     string   `json:"mods_path"`
	LauncherPath string   `json:"launcher_path"`
	AccessKey    string   `json:"access_key"`
	SecretKey    string   `json:"secret_key"`
	MinioHost    string   `json:"minio_host"`
}

func (manifest *Manifest) Write() error {
	byteArr, err := json.Marshal(manifest.Records)
	if err != nil {
		return err
	}
	ioutil.WriteFile(manifestFilename, byteArr, os.FileMode(0644))
	return nil
}

func (manifest *Manifest) Read() error {
	data, err := ioutil.ReadFile(manifestFilename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &manifest.Records)
	if err != nil {
		return err
	}
	return nil
}

func CreateConfigTemplate() error {
	data := []byte(`{
	"excluded": [],
	"mods_path": "",
	"launcher_path": "",
	"access_key": "",
	"secret_key": "",
	"minio_host": ""
}`)
	return ioutil.WriteFile("config.json", data, os.FileMode(0644))
}

func (config *Config) Read() error {
	data, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	return nil
}
