package ref

import "time"

func Bool(i bool) *bool {
	return &i
}

func Uint(i uint) *uint {
	return &i
}

func Uint8(i uint8) *uint8 {
	return &i
}

func Uint32(i uint32) *uint32 {
	return &i
}

func Uint64(i uint64) *uint64 {
	return &i
}

func Int(i int) *int {
	return &i
}

func Int8(i int8) *int8 {
	return &i
}

func Int32(i int32) *int32 {
	return &i
}

func Int64(i int64) *int64 {
	return &i
}

func String(i string) *string {
	return &i
}

func Time(i time.Time) *time.Time {
	return &i
}

func Duration(i time.Duration) *time.Duration {
	return &i
}

func Strings(ss []string) []*string {
	r := make([]*string, len(ss))
	for i := range ss {
		r[i] = &ss[i]
	}
	return r
}
