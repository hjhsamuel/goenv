package goenv

import (
	"github.com/pkg/errors"
	"os"
	"reflect"
	"strconv"
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
	default:
		if tag.Name == e.prefix {
			return nil
		}
		var val string
		if val = os.Getenv(tag.Name); val == "" {
			if v.IsZero() {
				val = tag.Default
			}
		}
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
			if val != "" {
				if value, err := strconv.Atoi(val); err != nil {
					return err
				} else {
					v.SetInt(int64(value))
				}
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if val != "" {
				if value, err := strconv.Atoi(val); err != nil {
					return err
				} else {
					v.SetUint(uint64(value))
				}
			}
		case reflect.Float32, reflect.Float64:
			if val != "" {
				if value, err := strconv.ParseFloat(val, 64); err != nil {
					return err
				} else {
					v.SetFloat(value)
				}
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

func NewEnvParser() *EnvParser {
	return &EnvParser{
		prefix:    EnvDefaultPrefix,
		splitChar: EnvDefaultSplitChar,
		tag:       EnvDefaultTagName,
	}
}
