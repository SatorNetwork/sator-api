package httpencoder

import (
	"context"
	"encoding/json"
	"net/http"
)

// Predefined http encoder content type
const (
	ContentTypeHeader = "Content-Type"
	ContentType       = "application/json; charset=utf-8"
)

type (
	// Response struct
	Response struct {
		Data interface{}            `json:"data,omitempty"`
		Meta map[string]interface{} `json:"meta,omitempty"`
	}

	// ListResponse struct
	ListResponse struct {
		Data interface{} `json:"data,omitempty"`
		Meta struct {
			TotalItems   int64 `json:"total_items,omitempty"`
			ItemsPerPage int64 `json:"items_per_page,omitempty"`
			Page         int64 `json:"page,omitempty"`
			Limit        int64 `json:"limit,omitempty"`
			Offset       int64 `json:"offset,omitempty"`
		} `json:"meta,omitempty"`
	}

	// BoolResultResponse struct
	BoolResultResponse struct {
		Result bool `json:"result"`
	}
)

// EncodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set(ContentTypeHeader, ContentType)

	if response == nil {
		w.WriteHeader(http.StatusCreated)
		return nil
	}

	switch response.(type) {
	case *Response, *BoolResultResponse, *ListResponse:
		return json.NewEncoder(w).Encode(response)
	}
	return json.NewEncoder(w).Encode(Response{Data: response})
}

// BoolResult response helper
func BoolResult(result bool) BoolResultResponse {
	return BoolResultResponse{Result: result}
}
