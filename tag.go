package goenv

import (
	"strings"
)

const (
	TagSplitChar = ";"

	TagValueChar  = ":"
	TagDefaultSig = "default"

	TagSkipSig     = "-"
	TagRequiredSig = "required"
)

type Tag struct {
	Name     string
	Default  string
	Required bool
}

func ParseTag(info string) (*Tag, error) {
	if info == "" {
		return nil, nil
	}
	vals := strings.Split(info, TagSplitChar)
	res := &Tag{}
	for _, val := range vals {
		switch val {
		case TagSkipSig:
			continue
		case TagRequiredSig:
			res.Required = true
		default:
			l := strings.Split(val, TagValueChar)
			if len(l) != 2 {
				res.Name = l[0]
			} else {
				switch l[0] {
				case TagDefaultSig:
					res.Default = l[1]
				}
			}
		}
	}
	return res, nil
}
