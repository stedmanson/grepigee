package apigee

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

func GetItemList(url string) ([]string, error) {
	var list = new([]string)

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

func GetDeployments(list []string, environment, path string) chan Deployment {
	var out = make(chan Deployment)
	var wg sync.WaitGroup

	for _, item := range list {
		wg.Add(1)
		go func(item string) {
			defer wg.Done()

			url := baseURL + "organizations/woolworths/environments/" + environment + path + item + "/deployments"
			body, err := Get(url)
			if err != nil {
				if errors.Is(err, ErrBadRequest) {
					return
				} else {
					fmt.Printf("Error fetching deployments for %s: %v\n", item, err)
					return
				}
			}

			var deployment Deployment
			if path == "/apis/" {
				deployment = new(ProxyDeployment)
			} else if path == "/sharedflows/" {
				deployment = new(SharedflowDeployment)
			}

			err = json.Unmarshal(body, &deployment)
			if err != nil {
				fmt.Printf("Error unmarshalling deployment for %s: %v\n", item, err)
				return
			}

			out <- deployment
		}(item)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
