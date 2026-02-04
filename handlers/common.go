package handlers

import "github.com/danielgtaylor/huma/v2"

type Pagination struct {
	Limit  int32 `query:"limit" minimum:"1" required:"true" default:"10"`
	Offset int32 `query:"offset" required:"true" default:"0"`
}

var AccessDeniedError = huma.Error403Forbidden("Access denied: record does not exist or you do not have the right to perform this operation")
