package apigee

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

func humanReadable(n string, asTime bool) (string, error) {
	number, err := strconv.ParseFloat(n, 64)
	if err != nil {
		return "", fmt.Errorf("error converting string to float64: %v", err)
	}

	intNumber := int64(number)

	p := message.NewPrinter(language.English)

	if asTime {
		duration := time.Duration(intNumber) * time.Millisecond // Assuming input is in milliseconds
		if duration.Hours() >= 1 {
			return fmt.Sprintf("%dh %dm %ds", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60), nil
		} else if duration.Minutes() >= 1 {
			return fmt.Sprintf("%dm %ds", int(duration.Minutes()), int(duration.Seconds())%60), nil
		} else if duration.Seconds() >= 1 {
			return fmt.Sprintf("%ds", int(duration.Seconds())), nil
		} else {
			return fmt.Sprintf("%dms", int(duration.Milliseconds())), nil
		}
	} else {
		return p.Sprintf("%d", intNumber), nil
	}
}
