package goenv

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

// common
const (
	TagSplitChar = ";"

	TagSliceSplitChar = ","
	TagMapSplitChar   = "|"

	TagName = "name"

	TagValueChar  = ":"
	TagDefaultSig = "default"

	TagSkipSig     = "-"
	TagRequiredSig = "required"
)

type Tag struct {
	Name     string
	Default  string
	Required bool
	Number   *NumberTag
}

// number
const (
	TagNumberLessThan    = "lt"
	TagNumberLessOrEqual = "lte"

	TagNumberGreaterThan    = "gt"
	TagNumberGreaterOrEqual = "gte"
)

type NumberTag struct {
	LessThan    *float64
	LessOrEqual *float64

	GreaterThan    *float64
	GreaterOrEqual *float64
}

func (t *NumberTag) set(key, value string) error {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	switch key {
	case TagNumberLessThan:
		t.LessThan = &v
	case TagNumberLessOrEqual:
		t.LessOrEqual = &v
	case TagNumberGreaterThan:
		t.GreaterThan = &v
	case TagNumberGreaterOrEqual:
		t.GreaterOrEqual = &v
	}
	return nil
}

func (t *NumberTag) check(val float64) error {
	if t.LessThan != nil && val >= *t.LessThan {
		return errors.Errorf("%f should be less than %f", val, *t.LessThan)
	}
	if t.LessOrEqual != nil && val > *t.LessOrEqual {
		return errors.Errorf("%f should be less or equal to %f", val, *t.LessOrEqual)
	}
	if t.GreaterThan != nil && val <= *t.GreaterThan {
		return errors.Errorf("%f should be greater than %f", val, *t.GreaterThan)
	}
	if t.GreaterOrEqual != nil && val < *t.GreaterOrEqual {
		return errors.Errorf("%f should be greater or equal to %f", val, *t.GreaterOrEqual)
	}
	return nil
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
			index := strings.Index(val, TagValueChar)
			if index == -1 {
				res.Name = val
			} else {
				tagVal := strings.TrimSpace(val[index+1:])
				switch val[:index] {
				case TagName:
					res.Name = tagVal
				case TagDefaultSig:
					res.Default = tagVal
				case TagNumberLessThan, TagNumberLessOrEqual, TagNumberGreaterThan, TagNumberGreaterOrEqual:
					if res.Number == nil {
						res.Number = &NumberTag{}
					}
					if err := res.Number.set(val[:index], tagVal); err != nil {
						return nil, err
					}
				}
			}

			//l := strings.Split(val, TagValueChar)
			//if len(l) < 2 {
			//	res.Name = l[0]
			//} else {
			//	tagVal := strings.TrimSpace(l[1])
			//	switch l[0] {
			//	case TagName:
			//		res.Name = tagVal
			//	case TagDefaultSig:
			//		res.Default = tagVal
			//	case TagNumberLessThan, TagNumberLessOrEqual, TagNumberGreaterThan, TagNumberGreaterOrEqual:
			//		if res.Number == nil {
			//			res.Number = &NumberTag{}
			//		}
			//		if err := res.Number.set(l[0], tagVal); err != nil {
			//			return nil, err
			//		}
			//	}
			//}
		}
	}
	return res, nil
}
