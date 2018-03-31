package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"minecraft_prelauncher/helpers"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	assetsPath := "assets/bibliocraft/textures/custompaintings"
	source := flag.String("source", "./", "a filepath to minio folder")
	destination := flag.String("destination", "./", "a filepath to destination jar file")
	flag.Parse()

	tempFilepaths := make([]string, 0)
	images, err := readImages(*source)
	if err != nil {
		fmt.Println("Could not read the images", err)
		os.Exit(1)
	}
	for _, img := range images {
		newFilename := fmt.Sprintf("%s%s", helpers.CalculateChecksum(img), ".png")
		newFilepath := path.Join(assetsPath, newFilename)
		tempFilepaths = append(tempFilepaths, newFilepath)
		cmd := exec.Command("cp", img, newFilepath)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Could not copy the file", err)
			os.Exit(1)
		}
		err := os.Remove(img)
		if err != nil {
			fmt.Println("Failed to delete the file", err)
			os.Exit(1)
		}
	}

	for _, fp := range tempFilepaths {
		cmd := exec.Command("zip", *destination, fp)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Could not zip the image", err)
			os.Exit(1)
		}
		err := os.Remove(fp)
		if err != nil {
			fmt.Println("Failed to delete the file", err)
			os.Exit(1)
		}
	}

}

func readImages(root string) ([]string, error) {
	var files []string

	fileInfos, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, info := range fileInfos {
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".png") {
			continue
		}
		files = append(files, filepath.Join(root, info.Name()))
	}
	return files, nil
}
