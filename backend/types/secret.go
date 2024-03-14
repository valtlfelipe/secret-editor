package types

type Secret struct {
	ARN    string `json:"arn"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type Result struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Result  Secret `json:"result"`
}

type ResultList struct {
	Success bool     `json:"success"`
	Error   string   `json:"error"`
	Result  []Secret `json:"result"`
}
