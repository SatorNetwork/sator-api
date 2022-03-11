package utils

import "net/http"

func IsStatusCodeSuccess(code int) bool {
	return code >= http.StatusOK && code < 300
}
