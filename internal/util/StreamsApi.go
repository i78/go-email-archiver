package util

import "io/fs"

func Map(ss []string, fn func(string) string) (result []string) {
	for _, s := range ss {
		result = append(result, fn(s))
	}
	return
}

func ForEach(ss []string, fn func(string)) {
	for _, s := range ss {
		fn(s)
	}
}

func Max(oo *[]fs.FileInfo, comparator func(o interface{}, p interface{}) bool) (result fs.FileInfo) {
	for _, item := range *oo {
		if result == nil || comparator(result, item) {
			result = item
		}
	}
	return
}
