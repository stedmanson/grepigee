package apigee

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	ErrBadRequest      = errors.New("bad request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrNotFound        = errors.New("not found")
	ErrInternalServer  = errors.New("internal server error")
	ErrUnknownResponse = errors.New("unknown response status")
)

func checkResponseStatus(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusInternalServerError:
		return ErrInternalServer
	default:
		return fmt.Errorf("received unexpected HTTP status: %d", resp.StatusCode)
	}
}

func Get(url string) ([]byte, error) {
	fmt.Println(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}

func DownloadBinaryContent(url string, outputPath string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp); err != nil {
		return err
	}

	// Create the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	// Write the response body to the output file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing to output file: %v", err)
	}

	return nil
}
