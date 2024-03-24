package apigee

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetSharedFlowList() ([]string, error) {
	return GetItemList(baseURL + "/sharedflows/")
}

func GetSharedflowDeployments(list []string, environment string) chan SharedflowDeployment {
	genericDeployments := GetDeployments(list, environment, "/sharedflows/")
	specificDeployments := make(chan SharedflowDeployment)

	go func() {
		for deployment := range genericDeployments {
			if sfd, ok := deployment.(*SharedflowDeployment); ok {
				specificDeployments <- *sfd
			}
		}
		close(specificDeployments)
	}()

	return specificDeployments
}

func DownloadSharedflowRevision(in chan SharedflowDeployment, environment string) {
	dirPath := filepath.Join(environment, "sharedflows")

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	for sharedflow := range in {
		if sharedflow.Environment != environment {
			continue
		}

		for _, deployment := range sharedflow.Revision {
			if deployment.State != "deployed" {
				continue
			}
			url := baseURL + "/sharedflows/" + sharedflow.Name + "/revisions/" + deployment.Name + "?format=bundle"

			folderName := fmt.Sprintf("%s-%s", sharedflow.Name, deployment.Name)
			outputPath := filepath.Join(dirPath, folderName+".zip")

			// Download the file
			if err := DownloadBinaryContent(url, outputPath); err != nil {
				fmt.Printf("Error downloading %s: %v\n", url, err)
				continue
			}

			// Unzip the file into a folder named after the sharedflow and revision
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
