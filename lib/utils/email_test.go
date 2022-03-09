package utils

import "testing"

func TestSanitizeEmail(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test@mail.dev", args{"test@mail.dev"}, "test@mail.dev", false},
		{"trim plus", args{"test+1@mail.dev"}, "test@mail.dev", false},
		{"trim spaces", args{" test@mail.dev  "}, "test@mail.dev", false},
		{"trim plus", args{"test+test@mail.dev"}, "test@mail.dev", false},
		{"trim dots", args{"tes.t@mail.dev"}, "test@mail.dev", false},
		{"trim dots 2", args{"t.e.s.t@mail.dev"}, "test@mail.dev", false},
		{"trim dash", args{"te-st@mail.dev"}, "test@mail.dev", false},
		{"full test", args{"tes.t+23@mail.dev"}, "test@mail.dev", false},
		{"wrong email", args{"tes.t+23.mail.dev"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SanitizeEmail(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizeEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SanitizeEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
