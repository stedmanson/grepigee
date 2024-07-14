package apigee

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func GetProxyList() ([]string, error) {
	return GetItemList(baseURL + "organizations/woolworths/apis/")
}

func StreamProxyDeployments(list []string, environment string) (chan ProxyDeployment, chan UndeployedEntity) {
	genericDeployments := GetDeployments(list, environment, "/apis/")
	specificDeployments := make(chan ProxyDeployment)
	undeployedEntities := make(chan UndeployedEntity)

	go func() {
		defer close(specificDeployments)
		defer close(undeployedEntities)

		for deployment := range genericDeployments {
			switch d := deployment.(type) {
			case *UndeployedEntity:
				undeployedEntities <- *d
			case *ProxyDeployment:
				specificDeployments <- *d
			default:
				fmt.Printf("Unknown deployment type: %T\n", deployment)
			}
		}
	}()

	return specificDeployments, undeployedEntities
}

func DownloadProxyRevision(in chan ProxyDeployment, environment string) chan struct{} {
	done := make(chan struct{})
	dirPath := filepath.Join(environment, "proxies")

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		close(done)
		return done
	}

	var wg sync.WaitGroup

	go func() {
		for proxy := range in {
			if len(proxy.Revision) == 0 {
				continue
			}
			for _, deployment := range proxy.Revision {
				if deployment.State != "deployed" {
					continue
				}

				wg.Add(1)

				go func(proxyName, deploymentName string) {
					defer wg.Done()
					url := baseURL + "organizations/woolworths/apis/" + proxyName + "/revisions/" + deploymentName + "?format=bundle"

					folderName := fmt.Sprintf("%s-%s", proxyName, deploymentName)
					outputPath := filepath.Join(dirPath, folderName+".zip")

					// Download the file
					if err := DownloadBinaryContent(url, outputPath); err != nil {
						fmt.Printf("Error downloading %s: %v\n", url, err)
						return
					}

					// Unzip the file into a folder named after the proxy and revision
					unzipPath := filepath.Join(dirPath, folderName)
					if err := os.MkdirAll(unzipPath, os.ModePerm); err != nil {
						fmt.Printf("Error creating directory for unzipping: %v\n", err)
						return
					}
					if err := Unzip(outputPath, unzipPath); err != nil {
						fmt.Printf("Error unzipping %s: %v\n", outputPath, err)
						return
					}
				}(proxy.Name, deployment.Name)
			}
		}

		wg.Wait()
		close(done)
	}()

	return done
}
