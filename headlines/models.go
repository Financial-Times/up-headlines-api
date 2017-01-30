package headlines

type headlineOutput struct {
	UUID       string `json:"uuid"`
	Title      string `json:"title"`
	Standfirst string `json:"standfirst"`
}

type HeadlineInput struct {
	UUIDs []string `json:"uuids"`
}
