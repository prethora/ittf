package pcre2regexp

import "github.com/Jemmic/go-pcre2"

type Regexp struct {
	matcher *pcre2.Matcher
}

func Compile(pattern string) (*Regexp, error) {
	regexp, err := pcre2.Compile(pattern, 0)
	if err != nil {
		return nil, err
	}
	matcher := regexp.NewMatcher()
	ret := Regexp{matcher: matcher}

	return &ret, nil
}

func (r *Regexp) FindAllStringSubmatchIndex(s string) [][]int {
	ret := [][]int{}
	offset := 0

	for {
		if r.matcher.MatchString(s, 0) {
			single := []int{}
			gn := r.matcher.Groups()
			var afterIndex int
			for i := 0; i <= gn; i++ {
				indices := r.matcher.GroupIndices(i)
				single = append(single, indices[0]+offset, indices[1]+offset)
				if i == 0 {
					afterIndex = indices[1]
				}
			}
			ret = append(ret, single)
			s = s[afterIndex:]
			offset += afterIndex
		} else {
			break
		}
	}

	if len(ret) == 0 {
		return nil
	} else {
		return ret
	}
}
