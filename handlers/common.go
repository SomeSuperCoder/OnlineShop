package handlers

type Pagination struct {
	Limit  int32 `query:"limit" minimum:"1" required:"true" default:"10"`
	Offset int32 `query:"offset" required:"true" default:"0"`
}
