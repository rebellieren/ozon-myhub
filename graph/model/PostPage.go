package model

type PostPage struct {
	Posts       []*Post `json:"posts"`
	TotalCount  int32   `json:"totalCount"`
	HasNextPage bool    `json:"hasNextPage"`
}
