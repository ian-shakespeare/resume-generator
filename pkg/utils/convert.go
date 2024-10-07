package utils

func StringToRef(s string) *string {
	var ref *string
	if len(s) >= 1 {
		ref = &s
	}
	return ref
}
