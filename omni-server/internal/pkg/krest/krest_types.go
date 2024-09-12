package krest

// Resources
type ResourceLinks struct {
	Self string `json:"self"`
}

type ResourceQuery struct {
	Expand []string `json:"expand"`
}

type ResourceResponse[T any] struct {
	Links      ResourceLinks     `json:"@links"`
	Query      ResourceQuery     `json:"@query"`
	Expandable map[string]string `json:"@expandable"`
	Resource   T
}

// Collections
type CollectionLinks struct {
	Self     string `json:"self"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type CollectionQuery struct {
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
	Expand []string `json:"expand"`
}

type CollectionResponse[T any] struct {
	Links      CollectionLinks       `json:"@links"`
	Query      CollectionQuery       `json:"@query"`
	Expandable map[string]string     `json:"@expandable"`
	Count      int                   `json:"count"`
	Total      int                   `json:"total"`
	Results    []ResourceResponse[T] `json:"results"`
}
