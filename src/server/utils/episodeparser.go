package utils

import (
	"errors"
	"strconv"
	"strings"
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

func getTokens(s string) (tokens []string, err error) {
	ignore := map[rune]int{
		'(': 0,
		')': 1,
		'[': 2,
		']': 3,
		'【': 4,
		'】': 5,
	}

	stack := []rune{}
	var token string

	addToken := func() {
		if len(token) > 0 {
			tokens = append(tokens, token)
			token = ""
		}
	}

	for _, r := range s {
		if p, ok := ignore[r]; ok {
			addToken()
			if p%2 == 0 {
				stack = append(stack, r)
			} else {
				if len(stack) == 0 || ignore[stack[len(stack)-1]]+1 != p {
					return nil, errors.New("")
				}
				stack = stack[:len(stack)-1]
				addToken()
			}
		} else {
			if r == ' ' && len(stack) == 0 {
				addToken()
				continue
			}
			token += string(r)
		}
	}
	return
}

func getNum(s string) (int, error) {
	ret := 0
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return -1, errors.New("not a digit")
		}
		ret = ret*10 + int(c-rune('0'))
	}
	return ret, nil
}

func chinese2Num(s string) (int, error) {
	m := map[rune]int{
		'零': 0,
		'一': 1,
		'二': 2,
		'三': 3,
		'四': 4,
		'五': 5,
		'六': 6,
		'七': 7,
		'八': 8,
		'九': 9,
		'十': 10,
		'百': 100,
		'千': 1000,
		'万': 10000,
		'亿': 100000,
	}
	ret := 0
	tmp := 0
	lastN := -1
	for _, c := range s {
		if n, ok := m[c]; !ok {
			return -1, nil
		} else {
			if tmp == 0 {
				if c == '零' {
					continue
				}
				tmp = n
				lastN = n
				continue
			}
			if n < lastN {
				ret += tmp
				tmp = n
				continue
			}
			tmp *= n
			lastN = n
		}
	}
	ret += tmp
	return ret, nil
}

func ParseEpisode2(name string) (int, error) {
	tokens, err := getTokens(name)
	if err != nil {
		return -1, err
	}
	for _, t := range tokens {
		n, err := getNum(t)
		if err == nil {
			return n, nil
		}
		if strings.HasPrefix(t, "第") && strings.HasSuffix(t, "集") {
			t = t[len("第"):]
			t = t[:len(t)-len("集")]
			n, err := getNum(t)
			if err == nil {
				return n, nil
			}
			n, err = chinese2Num(t)
			if err == nil {
				return n, nil
			}
		}
	}
	return -1, nil
}
