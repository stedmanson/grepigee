package apigee

import (
	"os"
)

var baseURL string
var token string

func init() {
	baseURL = "https://api.enterprise.apigee.com/v1/organizations/woolworths"
	token = os.Getenv("APIGEE_BEARER_TOKEN")
}
