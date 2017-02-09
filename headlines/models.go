package headlines

type headlineOutput struct {
	UUID       string `json:"uuid"`
	Title      string `json:"title"`
	Standfirst string `json:"standfirst"`
}

type HeadlineInput struct {
	UUIDs []string `json:"uuids,omitempty"`
}

type List struct {
	ID               string     `json:"id,omitempty"`
	Title            string     `json:"title,omitempty"`
	APIurl           string     `json:"apiUrl,omitempty"`
	ListType         string     `json:"listType,omitempty"`
	Items            []ListItem `json:"items,omitempty"`
	LayoutHint       string     `json:"layoutHint,omitempty"`
	PublishReference string     `json:"publishReference,omitempty"`
	LastModified     string     `json:"lastModified,omitempty"`
}

type ListItem struct {
	ID     string `json:"id,omitempty"`
	APIurl string `json:"apiUrl,omitempty"`
}
