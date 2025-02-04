package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"encoding/json"
    	"os"
)

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func detectImages(m map[string]string) {
	fmt.Println("\n--- Searching images in K8S manifests ---")
	// Populate keys (filenames)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	var filesWithImg []string
	var linesWithImg []string
	var imageList []string

	for _, k := range keys { // for K8S manifest file
		linesWithImg = []string{}
		scanner := bufio.NewScanner(strings.NewReader(m[k])) // read K8S manifest file

		for scanner.Scan() {
			cleanLine := strings.TrimSpace(scanner.Text())

			if !(strings.HasPrefix(cleanLine, "#")) {
				filesWithImg = append(filesWithImg, k)
				re := regexp.MustCompile(`image:.+`)
				linesWithImg = re.FindAllString(cleanLine, -1)

				if len(linesWithImg) > 0 {
					fmt.Printf("\nImage found in %s...\n", k)
				}

				for _, i := range linesWithImg {
					image := strings.TrimPrefix(i, "image:")
					image = strings.TrimSpace(image)
					image = strings.Trim(image, "\"")
					fmt.Println(image)
					imageList = append(imageList, image)
				}

			}
		}

	}

	uniqueImageList := unique(imageList)

	fmt.Println("\n--- List of unique Docker images in the Helm Chart ---\n")
	for _, i := range uniqueImageList {
		fmt.Println("\u2192", i)
	}

	imageDict := map[string][]string{"images": uniqueImageList}

	jsonFile, err := os.Create("images.json")
    	if err != nil {
		fmt.Println(err)
		return
    	}
    	defer jsonFile.Close()

    	encoder := json.NewEncoder(jsonFile)
    	err = encoder.Encode(imageDict)
    	if err != nil {
		fmt.Println(err)
		return
    	}

}
