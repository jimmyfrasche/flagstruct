package flagstruct

import (
	"go/ast"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type flg struct {
	defaultVal interface{}
	name, dscr string
	idx        []int
}

//needed to compare against to disambiguate the int64 kind
var durationType = reflect.ValueOf(time.Nanosecond).Type()

//used for int parsing
var intSize = 64

func init() {
	switch runtime.GOARCH {
	case "386", "arm":
		intSize = 32
	}
}

func parseType(t reflect.Type) ([]flg, error) {
	var flags []flg
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		//we only care about exported fields that are
		//either structs
		//or a valid kind that has a flag tag
		if !ast.IsExported(f.Name) {
			continue
		}
		k := f.Type.Kind()
		if k == reflect.Struct {
			sub, err := parseType(f.Type)
			if err != nil {
				return nil, errw(err)
			}
			//copy sub's flags into ours, prefixing i to their indicies
			for _, s := range sub {
				flags = append(flags, flg{
					defaultVal: s.defaultVal,
					name:       s.name,
					dscr:       s.dscr,
					idx:        append([]int{i}, s.idx...),
				})
			}
		} else if tag := f.Tag.Get("flag"); validKind(k) && tag != "" {
			name, dv, dscr := parseFlagKey(tag)

			//if no name given, rewrite the field name
			if name == "" {
				name = rewriteName(f.Name)
			}

			//parse the default value
			var defaultVal interface{} = dv
			s := strings.TrimSpace(dv) //for nonstrings we don't want extra ws
			if f.Type == durationType {
				if s == "" {
					defaultVal = time.Duration(0)
				} else {
					d, err := time.ParseDuration(s)
					if err != nil {
						return nil, errw(err)
					}
					defaultVal = d
				}
			} else {
				switch k {
				case reflect.String:
					//handled when defaultVal defined
				case reflect.Bool:
					if s == "" {
						defaultVal = false
						break
					}
					b, err := strconv.ParseBool(s)
					if err != nil {
						return nil, errw(err)
					}
					defaultVal = b
				case reflect.Int:
					if s == "" {
						defaultVal = int(0)
						break
					}
					i, err := strconv.ParseInt(s, 10, intSize)
					if err != nil {
						return nil, errw(err)
					}
					defaultVal = int(i)
				case reflect.Int64:
					if s == "" {
						defaultVal = int64(0)
						break
					}
					i, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return nil, errw(err)
					}
					defaultVal = i
				case reflect.Uint:
					if s == "" {
						defaultVal = uint(0)
						break
					}
					u, err := strconv.ParseUint(s, 10, intSize)
					if err != nil {
						return nil, errw(err)
					}
					defaultVal = uint(u)
				case reflect.Uint64:
					if s == "" {
						defaultVal = uint64(0)
						break
					}
					u, err := strconv.ParseUint(s, 10, 64)
					if err != nil {
						return nil, errw(err)
					}
					defaultVal = u
				case reflect.Float64:
					if s == "" {
						defaultVal = float64(0)
						break
					}
					f, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return nil, errw(err)
					}
					defaultVal = f
				}
			}

			//build our flag
			flags = append(flags, flg{
				defaultVal: defaultVal,
				name:       name,
				dscr:       dscr,
				idx:        f.Index,
			})
		}
	}

	return flags, nil
}

func validKind(v reflect.Kind) bool {
	switch v {
	case reflect.Bool,
		reflect.Int,
		reflect.Int64, //also covers time.Duration
		reflect.Uint,
		reflect.Uint64,
		reflect.Float64,
		reflect.String:
		return true
	}
	return false
}

func parseFlagKey(s string) (name, dv, dscr string) {
	sc := ","
	var sr rune
	for i, r := range s {
		if i == 0 {
			sr = r
		}
		if i == 1 {
			if r == ':' {
				sc = string([]rune{sr})
			}
			break
		}
	}
	v := strings.SplitN(s, sc, 3)
	ln := len(v)
	if ln > 0 {
		name = strings.TrimSpace(v[0])
	}
	if ln > 1 {
		dv = v[1]
	}
	if ln > 2 {
		dscr = v[2]
	}
	return
}

func rewriteName(nm string) string {
	return strings.ToLower(strings.Replace(nm, "_", "-", -1))
}
