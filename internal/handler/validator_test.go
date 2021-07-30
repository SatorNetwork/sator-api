package handler

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func Test_validate(t *testing.T) {
	data := url.Values{}
	data.Add("name", "John Doe")
	r, err := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	r.Form = data

	type args struct {
		req *http.Request
		rul map[string][]string
		msg map[string][]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no errors", args{r, map[string][]string{"name": []string{"required", "min:3"}}, nil}, false},
		{"email required", args{r, map[string][]string{"email": []string{"required", "email"}}, nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.args.req, tt.args.rul, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
