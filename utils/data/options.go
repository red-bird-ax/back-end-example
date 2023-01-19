package data

type Pagination struct {
	Offset int64
	Limit  int64
}

type Options struct {
	Pagination Pagination
	OrderBy    string
}