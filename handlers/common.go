package handlers

type Pagination struct {
	Limit  int32 `query:"limit"`
	Offset int32 `query:"offset"`
}
