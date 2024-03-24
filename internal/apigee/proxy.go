package apigee

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetProxyList() ([]string, error) {
	return GetItemList(baseURL + "/apis/")
}

func GetProxyDeployments(list []string, environment string) chan ProxyDeployment {
	genericDeployments := GetDeployments(list, environment, "/apis/")
	specificDeployments := make(chan ProxyDeployment)

	go func() {
		for deployment := range genericDeployments {
			if pd, ok := deployment.(*ProxyDeployment); ok {
				specificDeployments <- *pd
			}
		}
		close(specificDeployments)
	}()

	return specificDeployments
}

func DownloadProxyRevision(in chan ProxyDeployment, environment string) {

	dirPath := filepath.Join(environment, "proxies")

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	for proxy := range in {
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

}
