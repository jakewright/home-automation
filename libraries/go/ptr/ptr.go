package ptr

// Byte returns a *byte
func Byte(v byte) *byte { return &v }

// Int returns an *int
func Int(v int) *int { return &v }

// Int64 returns an *int64
func Int64(v int64) *int64 { return &v }

// Float64 returns a *float64
func Float64(v float64) *float64 { return &v }

// Bool returns a *bool
func Bool(v bool) *bool { return &v }

// String returns a *string
func String(v string) *string { return &v }
