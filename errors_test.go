package errors_test

import (
	"reflect"
	"testing"

	"github.com/tweithoener/errors"
)

// !!! most of the testing is done in the example_test.go file !!!

func TestAttributeFormatters(t *testing.T) {
	s := "abc"
	i := 123
	f := 123.456
	st := struct {
		A int
		B string
	}{789, "xyz"}
	format := "%s %d %7.3f %v"
	check := "abc 123 123.456 {789 xyz}"

	tests := []struct {
		name   string
		typ    string
		output interface{}
		check  interface{}
	}{
		{"Modf", "Mod", errors.Modf(format, s, i, f, st), errors.Mod(check)},
		{"Funcf", "Func", errors.Funcf(format, s, i, f, st), errors.Func(check)},
		{"Objf", "Obj", errors.Objf(format, s, i, f, st), errors.Obj(check)},
		{"Opf", "Op", errors.Opf(format, s, i, f, st), errors.Op(check)},
		{"Kindf", "Kind", errors.Kindf(format, s, i, f, st), errors.Kind(check)},
		{"Codef", "Code", errors.Codef(format, s, i, f, st), errors.Code(check)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.output != tt.check {
				t.Errorf("%s() = %v; expected %v", tt.name, tt.output, tt.check)
			}
			if typ := reflect.TypeOf(tt.output).Name(); typ != tt.typ {
				t.Errorf("reflect.TypeOf( %s() ) = %s ; expected %s", tt.name, typ, tt.typ)
			}
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	errs := make([]errors.Error, 5)
	for i := 1; i < 5; i++ {
		errs[i] = errors.E(i, errs[i-1])
	}
	for i := 4; i >= 1; i-- {
		if u := errs[i].Unwrap(); u != errs[i-1] {
			t.Errorf("%s.Unwrap() = %s; expected %s", errs[i].Error(), u.Error(), errs[i-1].Error())
		}
	}
}

func TestEWithNill(t *testing.T) {
	check := "hallo; test"
	if s := errors.E("hallo", nil, nil, nil, "test").Error(); check != s {
		t.Errorf("E(<with nil>) = %s; expected %s", s, check)
	}
}
