package flagstruct

import (
	"flag"
	"io"
	"reflect"
	"time"
)

//A Parser wraps a FlagSet with the information necessary to fill
//a struct, as passed to New or Extend, in with the parsed flags
type Parser struct {
	fs   *flag.FlagSet
	sig  []flg
	into reflect.Value
	tbl  map[string]reflect.Value
}

//New creates a Parser named name with the type in v.
func New(name string, v interface{}) (*Parser, error) {
	fs := flag.NewFlagSet(name, 0)
	return Extend(fs, v)
}

//Extend fs with the flags defined in v.
//
//When extending an existing flagset, Parse must be called on
//the returned parser even if the flagset has since been parsed.
//
//It is up to the caller to ensure there are no naming colisions.
func Extend(fs *flag.FlagSet, v interface{}) (*Parser, error) {
	//unwrap v
	V := reflect.ValueOf
	r := V(v)
	t := r.Type()
	if k := t.Kind(); k != reflect.Ptr {
		return nil, errf("v must be pointer to struct, got %s", k)
	}
	r = reflect.Indirect(r)
	t = r.Type()
	if k := t.Kind(); k != reflect.Struct {
		return nil, errf("v must be pointer to struct, got pointer to %s", k)
	}
	sig, err := parseType(t)
	if err != nil {
		return nil, err
	}

	//populate flagset
	tbl := map[string]reflect.Value{}
	for _, f := range sig {
		n := f.name
		switch v := f.defaultVal.(type) {
		case bool:
			tbl[n] = V(fs.Bool(n, v, f.dscr))
		case time.Duration:
			tbl[n] = V(fs.Duration(n, v, f.dscr))
		case float64:
			tbl[n] = V(fs.Float64(n, v, f.dscr))
		case string:
			tbl[n] = V(fs.String(n, v, f.dscr))
		case int:
			tbl[n] = V(fs.Int(n, v, f.dscr))
		case int64:
			tbl[n] = V(fs.Int64(n, v, f.dscr))
		case uint:
			tbl[n] = V(fs.Uint(n, v, f.dscr))
		case uint64:
			tbl[n] = V(fs.Uint64(n, v, f.dscr))
		}
	}

	return &Parser{fs, sig, r, tbl}, nil
}

//Parse parses flag definitions from the argument list,
//which should not include the command name.
//
//The return value will be ErrHelp if -help was set but not defined.
//
//This will populate the *struct passed to New.
//
//If a FlagSet was extended with Extend this will parse it but it is safe
//to call this regardless.
func (p *Parser) Parse(args []string) error {
	if !p.Parsed() {
		if err := p.fs.Parse(args); err != nil {
			return err
		}
	}

	for _, f := range p.sig {
		v := reflect.Indirect(p.tbl[f.name])
		p.into.FieldByIndex(f.idx).Set(v)
	}

	return nil
}

//PrintDefaults is like FlagSet.PrintDefaults to an io.Writer.
func (p *Parser) PrintDefaults(w io.Writer) {
	p.fs.SetOutput(w)
	p.fs.PrintDefaults()
}

//Proxy remaining helpful methods of flagset

//Arg is as defined on FlagSet.
func (p *Parser) Arg(i int) string {
	return p.fs.Arg(i)
}

//Args is as defined on FlagSet.
func (p *Parser) Args() []string {
	return p.fs.Args()
}

//Lookup is as defined on FlagSet.
func (p *Parser) Lookup(name string) *flag.Flag {
	return p.fs.Lookup(name)
}

//NArg is as defined on FlagSet.
func (p *Parser) NArg() int {
	return p.fs.NArg()
}

//NFlag is as defined on FlagSet.
func (p *Parser) NFlag() int {
	return p.fs.NFlag()
}

//Parsed is as defined on FlagSet.
func (p *Parser) Parsed() bool {
	return p.fs.Parsed()
}

//Visit is as defined on FlagSet.
func (p *Parser) Visit(fn func(*flag.Flag)) {
	p.fs.Visit(fn)
}

//VisitAll is as defined on FlagSet.
func (p *Parser) VisitAll(fn func(*flag.Flag)) {
	p.fs.VisitAll(fn)
}
