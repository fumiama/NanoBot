package nano

import (
	"flag"
	"reflect"
	"strings"
)

func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\r', '\n':
		return true
	}
	return false
}

type argType int

const (
	argNo argType = iota
	argSingle
	argQuoted
)

// ParseShell 将指令转换为指令参数.
// modified from https://github.com/mattn/go-shellwords
func ParseShell(s string) []string {
	var args []string
	buf := strings.Builder{}
	var escaped, doubleQuoted, singleQuoted, backQuote bool
	backtick := ""

	got := argNo

	for _, r := range s {
		if escaped {
			buf.WriteRune(r)
			escaped = false
			got = argSingle
			continue
		}

		if r == '\\' {
			if singleQuoted {
				buf.WriteRune(r)
			} else {
				escaped = true
			}
			continue
		}

		if isSpace(r) {
			if singleQuoted || doubleQuoted || backQuote {
				buf.WriteRune(r)
				backtick += string(r)
			} else if got != argNo {
				args = append(args, buf.String())
				buf.Reset()
				got = argNo
			}
			continue
		}

		switch r {
		case '`':
			if !singleQuoted && !doubleQuoted {
				backtick = ""
				backQuote = !backQuote
			}
		case '"':
			if !singleQuoted {
				if doubleQuoted {
					got = argQuoted
				}
				doubleQuoted = !doubleQuoted
			}
		case '\'':
			if !doubleQuoted {
				if singleQuoted {
					got = argSingle
				}
				singleQuoted = !singleQuoted
			}
		default:
			got = argSingle
			buf.WriteRune(r)
			if backQuote {
				backtick += string(r)
			}
		}
	}

	if got != argNo {
		args = append(args, buf.String())
	}

	return args
}

var (
	boolType    = reflect.TypeOf(false)
	intType     = reflect.TypeOf(0)
	stringType  = reflect.TypeOf("")
	float64Type = reflect.TypeOf(float64(0))
)

func registerFlag(t reflect.Type, v reflect.Value) *flag.FlagSet {
	v = v.Elem()
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Tag.Get("flag")
		if name == "" {
			continue
		}
		help := field.Tag.Get("help")
		switch field.Type {
		case boolType:
			fs.BoolVar(v.Field(i).Addr().Interface().(*bool), name, false, help)
		case intType:
			fs.IntVar(v.Field(i).Addr().Interface().(*int), name, 0, help)
		case stringType:
			fs.StringVar(v.Field(i).Addr().Interface().(*string), name, "", help)
		case float64Type:
			fs.Float64Var(v.Field(i).Addr().Interface().(*float64), name, 0, help)
		default:
			panic("unsupported type")
		}
	}
	return fs
}
