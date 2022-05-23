package utils

// BoolPointer converts a bool to a pointer
func BoolPointer(b bool) *bool {
	return &b
}

// StringPointer converts a string to a pointer
func StringPointer(s string) *string {
	return &s
}
