package utils

import (
	"strconv"
	"unicode"
)

func ParseEpisode(names []string) map[int]int {
	sortFunc := func(i int) map[int]int {
		name := names[i]
		if len(name) == 0 {
			return map[int]int{}
		}
		selected := make(map[int]bool)
		for j := range names {
			if len(names[j]) > 0 && names[j][0] == name[0] {
				selected[j] = true
			}
		}
		if len(selected) <= 1 {
			return map[int]int{}
		}
		j := 1
		for ; j < len(name); j++ {
			b := false
			for k := range selected {
				if len(names[k]) <= j {
					delete(selected, k)
					continue
				}
				if names[k][j] != name[j] {
					if !unicode.IsDigit(rune(name[j])) {
						delete(selected, k)
						continue
					}
					b = true
				}
			}
			if b {
				break
			}
		}
		for ; j > 0 && unicode.IsDigit(rune(name[j-1])); j-- {
		}
		ret := make(map[int]int)
		for k := range selected {
			if !unicode.IsDigit(rune(names[k][j])) {
				continue
			}
			e := j + 1
			for ; e < len(names[k]) && unicode.IsDigit(rune(names[k][e])); e++ {
			}
			ep, err := strconv.Atoi(names[k][j:e])
			if err != nil {
				continue
			}
			ret[ep] = k
		}
		return ret
	}
	var ret map[int]int
	for i := range names {
		retTmp := sortFunc(i)
		if len(retTmp) > len(ret) {
			ret = retTmp
		}
	}
	return ret
}
