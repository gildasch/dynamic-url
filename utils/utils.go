package utils

import "math/rand"

func OneOf(ss []string) string {
	if len(ss) == 0 {
		return ""
	}

	return ss[rand.Intn(len(ss)-1)]
}
