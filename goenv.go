package goenv

import (
	"encoding"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	EnvDefaultPrefix    = "ENV"
	EnvDefaultSplitChar = "_"
	EnvDefaultTagName   = "env"
)

type EnvParser struct {
	prefix    string
	splitChar string
	tag       string
}

func (e *EnvParser) SetPrefix(name string) {
	e.prefix = name
}

func (e *EnvParser) SetSplitChar(name string) {
	e.splitChar = name
}

func (e *EnvParser) SetTag(name string) {
	e.tag = name
}

func (e *EnvParser) Start(conf interface{}) error {
	t, v := reflect.TypeOf(conf), reflect.ValueOf(conf)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.Errorf("invalid type %s", t)
	}
	return e.parse(t, v, &Tag{Name: e.prefix})
}

func (e *EnvParser) parse(t reflect.Type, v reflect.Value, tag *Tag) error {
	switch t.Kind() {
	case reflect.Ptr:
		if !v.Elem().IsValid() {
			v.Set(reflect.New(t.Elem()))
		}
		return e.parse(t.Elem(), v.Elem(), tag)
	case reflect.Struct:
		if has, err := e.setUnmarshal(v, tag); err != nil {
			return err
		} else if !has {
			for i := 0; i < t.NumField(); i++ {
				tagInfo, err := ParseTag(t.Field(i).Tag.Get(e.tag))
				if err != nil {
					return err
				}
				if tagInfo == nil || tagInfo.Name == "" {
					continue
				}
				if tag != nil {
					tagInfo.Name = tag.Name + e.splitChar + tagInfo.Name
				}
				err = e.parse(t.Field(i).Type, v.Field(i), tagInfo)
				if err != nil {
					return err
				}
			}
		}
	case reflect.Slice:
		fmt.Println(t.Elem().String(), v.String())
		return e.slice(t, v, tag)
	default:
		if tag.Name == e.prefix {
			return nil
		}
		val := e.getValue(v, tag)
		if val == "" && v.IsZero() && tag.Required {
			return errors.Errorf("param %s is required", tag.Name)
		}
		switch t.Kind() {
		case reflect.String:
			if val != "" {
				v.SetString(val)
			}
		case reflect.Bool:
			if val != "" {
				if value, err := strconv.ParseBool(val); err != nil {
					return err
				} else {
					v.SetBool(value)
				}
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value, isSet, err := e.numberCheck(val, tag.Number); err != nil {
				return err
			} else if isSet {
				v.SetInt(int64(value))
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value, isSet, err := e.numberCheck(val, tag.Number); err != nil {
				return err
			} else if isSet {
				v.SetUint(uint64(value))
			}
		case reflect.Float32, reflect.Float64:
			if value, isSet, err := e.numberCheck(val, tag.Number); err != nil {
				return err
			} else if isSet {
				v.SetFloat(value)
			}
		case reflect.Complex64, reflect.Complex128:
			if val != "" {
				if value, err := strconv.ParseComplex(val, 128); err != nil {
					return err
				} else {
					v.SetComplex(value)
				}
			}
		}
	}
	return nil
}

func (e *EnvParser) getValue(v reflect.Value, tag *Tag) string {
	var val string
	if val = os.Getenv(tag.Name); val == "" {
		if v.IsZero() {
			val = tag.Default
		}
	}
	return val
}

func (e *EnvParser) slice(t reflect.Type, v reflect.Value, tag *Tag) error {
	value := e.getValue(v, tag)
	l := strings.Split(value, TagSliceSplitChar)

	te := t.Elem()
	if te.Kind() == reflect.Ptr {
		te = te.Elem()
	}

	if _, ok := reflect.New(te).Interface().(encoding.TextUnmarshaler); ok {
		elemType := v.Type().Elem()
		result := reflect.MakeSlice(reflect.SliceOf(elemType), len(l), len(l))
		for i, data := range l {
			indexElem := result.Index(i)
			if indexElem.Kind() == reflect.Ptr {
				indexElem = reflect.New(elemType.Elem())
			} else {
				indexElem = indexElem.Addr()
			}
			if f, ok := indexElem.Interface().(encoding.TextUnmarshaler); ok {
				if err := f.UnmarshalText([]byte(data)); err != nil {
					return err
				}
			}
			if indexElem.Kind() == reflect.Ptr {
				result.Index(i).Set(indexElem)
			}
		}
		v.Set(result)
	} else {
		fmt.Println("system type: ", t.Elem(), t)
		result := reflect.MakeSlice(t, len(l), len(l))
		for i, data := range l {
			if err := e.parse(t.Elem(), result.Index(i), &Tag{Default: data}); err != nil {
				return err
			}
		}
		v.Set(result)
	}
	return nil
}

func (e *EnvParser) setUnmarshal(v reflect.Value, tag *Tag) (bool, error) {
	var hasUnmarshal = false
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
	} else if v.CanAddr() {
		v = v.Addr()
	}

	if f, ok := v.Interface().(encoding.TextUnmarshaler); ok {
		hasUnmarshal = true
		value := e.getValue(v, tag)
		if err := f.UnmarshalText([]byte(value)); err != nil {
			return hasUnmarshal, err
		}
	}
	return hasUnmarshal, nil
}

func (e *EnvParser) numberCheck(val string, tag *NumberTag) (float64, bool, error) {
	var (
		value   float64 = 0
		needSet         = false
		err     error
	)
	if val != "" {
		needSet = true
		value, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, false, err
		}
	}
	if tag != nil {
		if err = tag.check(value); err != nil {
			return 0, false, err
		}
	}
	return value, needSet, nil
}

func NewEnvParser(ops ...Options) *EnvParser {
	e := &EnvParser{
		prefix:    EnvDefaultPrefix,
		splitChar: EnvDefaultSplitChar,
		tag:       EnvDefaultTagName,
	}
	for _, op := range ops {
		op(e)
	}
	return e
}
