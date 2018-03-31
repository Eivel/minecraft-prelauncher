package main

import (
	"flag"
	"fmt"
	"minecraft_prelauncher/helpers"
	settings "minecraft_prelauncher/settings"
	"os"
	"path"
	"strings"
)

func main() {
	filepath := flag.String("filepath", "./", "a filepath to mods folder")
	flag.Parse()

	modPaths, err := helpers.ReadModPaths(*filepath)
	if err != nil {
		fmt.Println("Failed to read mods directory", err)
		os.Exit(1)
	}

	var manifest settings.Manifest

	for _, el := range modPaths {
		if !strings.HasSuffix(el, ".jar") {
			continue
		}
		manifest.Records = append(manifest.Records, settings.Record{Name: path.Base(el), Checksum: helpers.CalculateChecksum(el)})
	}

	err = manifest.Write()
	if err != nil {
		fmt.Println("Failed to save manifest.json", err)
		os.Exit(1)
	}
}
