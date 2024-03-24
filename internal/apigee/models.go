package apigee

type Deployment interface{}
type ProxyDeployment struct {
	Environment  string `json:"environment"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Revision     []struct {
		Configuration struct {
			BasePath string        `json:"basePath"`
			Steps    []interface{} `json:"steps"`
		} `json:"configuration"`
		Name   string        `json:"name"`
		Server []interface{} `json:"server"`
		State  string        `json:"state"`
	} `json:"revision"`
}

type SharedflowDeployment struct {
	Environment  string `json:"environment"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Revision     []struct {
		Configuration struct {
			BasePath      string        `json:"basePath"`
			ConfigVersion string        `json:"configVersion"`
			Steps         []interface{} `json:"steps"`
		} `json:"configuration"`
		Name   string `json:"name"`
		Server []struct {
			Pod struct {
				Name   string `json:"name"`
				Region string `json:"region"`
			} `json:"pod"`
			Status string   `json:"status"`
			Type   []string `json:"type"`
			UUID   string   `json:"uUID"`
		} `json:"server"`
		State string `json:"state"`
	} `json:"revision"`
}
