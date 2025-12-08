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
	ImagePath   string  `json:"imagePath"`
	Topics      []Topic `json:"topics,omitzero"`
	ID          int     `json:"id"`
	TopicCount  int     `json:"topicsCount,omitzero"`
}

type Pagination struct {
	Page       int `json:"page"`
	TotalPages int `json:"totalPages"`
	TotalItems int `json:"totalItems"`
	NextPage   int `json:"nextPage"`
	PrevPage   int `json:"prevPage"`
}

type Logo struct {
	URL    string `json:"url"`
	ID     int    `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Topic struct {
	Title     string `json:"title"`
	CreatedAt string `json:"createdAt"`
	ID        int    `json:"id"`
}
