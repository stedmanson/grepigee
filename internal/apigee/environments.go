package apigee

func GetEnvironments() ([]string, error) {

	return GetItemList(baseURL + "organizations/woolworths/environments")
}
