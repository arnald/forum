package domain

type CategoryData struct {
	Data struct {
		Categories []Category `json:"categories"`
	} `json:"data"`
}

type Category struct {
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Slug        string  `json:"slug,omitzero"`
	Description string  `json:"description,omitzero"`
	Topics      []Topic `json:"topics,omitzero"`
	Logo        Logo    `json:"logo"`
	ID          int     `json:"id"`
	TopicCount  int     `json:"topicCount"`
}

type Logo struct {
	URL    string `json:"url"`
	ID     int    `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Topic struct {
	Title string `json:"title"`
	ID    int    `json:"id"`
}
