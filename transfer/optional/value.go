package optional

func Int(i int) *int {
	return &i
}

func Float(f float64) *float64 {
	return &f
}

func String(s string) *string {
	return &s
}

func Bool(b bool) *bool {
	return &b
}
