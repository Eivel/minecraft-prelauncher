package main

import (
	"fmt"
	helpers "minecraft_prelauncher/helpers"
	"minecraft_prelauncher/minio_api"
	settings "minecraft_prelauncher/settings"
	"os"
	"os/exec"
	"path"
	"strings"
)

const defaultBucket = "mods"

func main() {
	helpers.GenerateConfig()

	var config settings.Config
	err := config.Read()
	if err != nil {
		fmt.Println("Failed to read config file", err)
		os.Exit(1)
	}

	if !helpers.FileExists(config.LauncherPath) {
		fmt.Println("The launcher directory may not be correct. Please, check the path in config.json")
		os.Exit(1)
	}

	modPaths, err := helpers.ReadModPaths(config.ModsPath)
	if err != nil {
		fmt.Println("Failed to read mods directory", err)
		os.Exit(1)
	}

	if !helpers.ModPathsCorrect(modPaths) {
		fmt.Println("The mods directory may not be correct. Please, check the path in config.json")
		os.Exit(1)
	}

	minioClient, err := minio_api.InitializeClient(config.AccessKey, config.SecretKey, config.MinioHost)
	if err != nil {
		fmt.Println("Failed to initialize minio client", err)
		os.Exit(1)
	}

	maintenance, err := helpers.IsMaintenance(minioClient, defaultBucket)
	if err != nil {
		fmt.Println("Failed to fetch server status", err)
		os.Exit(1)
	}
	if maintenance {
		fmt.Println("Server is currently in maintenance mode, try again later")
		os.Exit(0)
	}

	err = minio_api.DownloadFile(minioClient, defaultBucket, "manifest.json", "manifest.json")
	if err != nil {
		fmt.Println("Failed to download remote manifest file", err)
		os.Exit(1)
	}

	var remoteManifest settings.Manifest

	err = remoteManifest.Read()
	if err != nil {
		fmt.Println("Failed to read local manifest file", err)
		os.Exit(1)
	}

	checksums := make([]string, 0)

	for _, mp := range modPaths {
		localChecksum := helpers.CalculateChecksum(mp)
		if helpers.IsFileDeprecated(remoteManifest, mp, localChecksum, config.Excluded) {
			err := os.Remove(mp)
			if err != nil {
				fmt.Println("Failed to delete a mod", err)
				os.Exit(1)
			}
		} else {
			checksums = append(checksums, localChecksum)
		}
	}

	for _, mod := range remoteManifest.Records {
		if helpers.IsNewRemoteFile(mod, checksums) {
			err := minio_api.DownloadFile(minioClient, defaultBucket, path.Join(config.ModsPath, mod.Name), mod.Name)
			if err != nil {
				fmt.Println("Failed to download a mod", err)
				os.Exit(1)
			}
		}
	}

	err = os.Remove("manifest.json")
	if err != nil {
		fmt.Println("Could not remove manifest.json file", err)
		os.Exit(1)
	}

	var cmd *exec.Cmd
	if strings.HasSuffix(config.LauncherPath, ".jar") {
		cmd = exec.Command("java", "-jar", config.LauncherPath)
	} else {
		cmd = exec.Command(config.LauncherPath)
	}
	err = cmd.Run()
	if err != nil {
		fmt.Println("Could not start the launcher. Please check the path in config.json", err)
		os.Exit(1)
	}
}
