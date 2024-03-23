package apigee

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func GetProxyList() ([]string, error) {
	var list = new([]string)
	url := baseURL + "/apis/"

	body, err := Get(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, list)
	if err != nil {
		return nil, err
	}

	return *list, nil
}

func GetProxyDeployments(list []string, environment string) chan ProxyDeployment {
	var out = make(chan ProxyDeployment) // Initialize the channel
	var wg sync.WaitGroup

	for _, proxy := range list {
		wg.Add(1)
		go func(proxy string) {
			defer wg.Done()

			url := baseURL + "/environments/" + environment + "/apis/" + proxy + "/deployments"

			body, err := Get(url)
			if err != nil {
				if errors.Is(err, ErrBadRequest) {
					return
				} else {
					fmt.Println(url, err)
				}
			}

			var pd = new(ProxyDeployment)

			err = json.Unmarshal(body, pd)
			if err != nil {
				fmt.Println(err)
				return
			}

			out <- *pd

		}(proxy)
	}

	go func() {
		wg.Wait()  // Wait for all goroutines to finish
		close(out) // Close the channel
	}()

	return out
}

func DownloadProxyRevision(in chan ProxyDeployment, environment string) {
	var validate []string

	dirPath := filepath.Join(environment, "proxies")

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	for proxy := range in {
		if len(proxy.Revision) > 1 {
			validate = append(validate, proxy.Name)
		}
		for _, deployment := range proxy.Revision {
			if deployment.State != "deployed" {
				continue
			}
			url := baseURL + "/apis/" + proxy.Name + "/revisions/" + deployment.Name + "?format=bundle"

			folderName := fmt.Sprintf("%s-%s", proxy.Name, deployment.Name)
			outputPath := filepath.Join(dirPath, folderName+".zip")

			// Download the file
			if err := DownloadBinaryContent(url, outputPath); err != nil {
				fmt.Printf("Error downloading %s: %v\n", url, err)
				continue
			}

			// Unzip the file into a folder named after the proxy and revision
			unzipPath := filepath.Join(dirPath, folderName)
			if err := os.MkdirAll(unzipPath, os.ModePerm); err != nil {
				fmt.Printf("Error creating directory for unzipping: %v\n", err)
				continue
			}
			if err := Unzip(outputPath, unzipPath); err != nil {
				fmt.Printf("Error unzipping %s: %v\n", outputPath, err)
				continue
			}

			// Delete the zip file
			if err := os.Remove(outputPath); err != nil {
				fmt.Printf("Error deleting %s: %v\n", outputPath, err)
			}
		}
	}

	fmt.Println("Validate:", validate)
}
