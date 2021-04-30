// Package errors provides error creation functions for descriptive, wrapped errors.
//
// Each error consists of set of error attributes describing the module and funtion
// in which the error occurred, the operation that failed, the object on which this
// operation was performed, an error code and a cause for this error (a wrapped error),
// plus an arbitratry number of details on the error.
//
// The E() function, which creates a new Error takes an arbitrary number of these
// attributes. An error with no attributes can be created as well as a fully specified error
// with all attributes.
//
// This package provides a template mechanism for error creation. An error
// template is created from a list of error attributes (using the T() function).
// The error attributes passed to the T() function will be common to all errors
// created from this template. When creating an error from
// a template (using the template's E() method) the error attributes of
// the template can be amended or overwritten with an arbitrary number of
// error attributes.
//
// Errors can be wrapped. Simply add an error to the call of E(). It can then be
// obtained via the Error.Unwrap() method or the Err field of the Error.
//
// Best Pratices:
//
// Create a file errors.go in your package and create an error template for
// your entire package there.
//   var etpl = errors.T(errors.Mod("<your package name>"))
// Then create templates for each function in your package that creates errors
//   var etReadConfig = etpl.T(errors.Func("read configuration"))
// Also add all operations, error kinds and error codes that you will need in your package to
// the errors.go file:
//   const (
//     eSendRequest = errors.Op("send http request")
//     eIllegalItem = errors.Kind("illegal config item id"))
//   )
// This keeps all strings organized in one place and allows to write clean looking error code:
//   if err != nil {
//     return nil, etReadConig.E(eIllegalItem, item, err)
//   }
// The error that is returned here, containes the package and function name (from the template) and
// also has information about the kind of problem (an illegal config item was detected), the name of
// the config item and wraps the cause for the error.
//
// All this information can be retreived from the fields of the resulting Error. If you need
// identifyable errors so that functions receiving the error can react differently on different kinds
// of errors, consider adding an error code to the error.
//   const (
//     EIllegalConfigItem errors.Code = 1000 + iota
//     EIllegalConfigValue errors.Code
//   )
//   // ...
//   return errors.E(EIllegalConfigValue, ...)
//
// There are usage examples in the documentation of the most important functions of this package.
// They can also be found in the examples_test.go file.
package errors

import (
	"fmt"
	"strings"
)

// Error holds error information. It is returned by the E() function.
// All error information included in the call to E() will be added to
// the respective fields of this struct
type Error struct {
	Module    Mod    // The module/package in which the error occured
	Function  Func   // The function in which the error occurred
	Kind      Kind   // The kind of error
	Operation Op     // The operation (in the Function) that caused the error
	Object    Obj    // The object on which operations were performed
	Code      Code   // An error code identifying the error
	Details   string // Further details that were passed into E()
	Err       error  // Wrapped error: the cause of this error
}

// E create a new error from a list of provided error attributes. Attributes of types
// Mod, Func, Obj, Op, Kind, and Code are stored in the respective fields
// of the returned Error. If an attribute type appears more then once the latter
// one takes precedence.
// If an attribute implementing the error interface the new Error
// wraps the provided error.
// String representations of all attributes of other types are added to the Details field.
// These details are separated by semicolons.
func E(args ...interface{}) Error {
	return e(Error{}, args...)
}

// Template is a template for Error creation. It stores
// common error attributes such as Mod, Func, or Obj and
// copies them into any Error that is created from this template.
//
// This is helpful when a piece of code returns similar errors in multiple
// places. The template can be created once. When an error is created,
// it is templated and only error attributes specific to the concrete error
// need to be passed to the E() function.
type Template Error

// T create a new Template from the provided list of error attributes.
//
// See errors.E() for details on the attribute list.
func T(args ...interface{}) Template {
	return Template(e(Error{}, args...))
}

// T creates a new Template from this Template ammending or overwriting this
// template's attributes with the provided list of Error attribnnutes.
//
// See errors.E() for details on the attribute list.
func (tpl Template) T(args ...interface{}) Template {
	return Template(e(Error(tpl), args...))
}

// E creates a new Error from this Template ammending or overwriting this
// template's attributes with the provided list of Error attribnnutes.
//
// See errors.E() for details on the attribute list.
func (tpl Template) E(args ...interface{}) Error {
	return e(Error(tpl), args...)
}

// Kind is an error attribute that describes the kind of error. E.g. write error, read error
type Kind string

// Kindf is a shortcut to Kind(fmt.Sprintf(...)). It produces a Kind with the given content.
func Kindf(f string, args ...interface{}) Kind {
	return Kind(fmt.Sprintf(f, args...))
}

// Code  is an error attribute that represents an error code. Codes are used to identify errors.
type Code string

// Codef is a shortcut to Code(fmt.Sprintf(...)). It produces a Code with the given content.
func Codef(f string, args ...interface{}) Code {
	return Code(fmt.Sprintf(f, args...))
}

// Mod is an error attribute that describes the program module in which an error happened.
type Mod string

// Modf is a shortcut to Mod(fmt.Sprintf(...)). It produces a Mod with the given content.
func Modf(f string, args ...interface{}) Mod {
	return Mod(fmt.Sprintf(f, args...))
}

// Func is an error attribute that describes the function in which an error occured.
// I.e. the calling function (callee) which checks the return value of a called function and
// composed the error. The called function can be stored as an Op (Operation).
type Func string

// Funcf is a shortcut to Func(fmt.Sprintf(...)). It produces a Func with the given content.
func Funcf(f string, args ...interface{}) Func {
	return Func(fmt.Sprintf(f, args...))
}

// Op is an error attribute that describes the operation during which an error occured.
// This is the call that failed inside the current function (which can be specified using Func).
type Op string

// Opf is a shortcut to Op(fmt.Sprintf(...)). It produces an Op with the given content.
func Opf(f string, args ...interface{}) Op {
	return Op(fmt.Sprintf(f, args...))
}

// Obj is an error attribute that describes the object on which the operation that led to an error was performed.
type Obj string

// Objf is a shortcut to Obj(fmt.Sprintf(...)). It produces an Obj with the given content.
func Objf(f string, args ...interface{}) Obj {
	return Obj(fmt.Sprintf(f, args...))
}

func e(err Error, args ...interface{}) Error {
	for _, arg := range args {
		if arg == nil {
			continue
		}
		switch x := arg.(type) {
		case Mod:
			err.Module = x
		case Func:
			err.Function = x
		case Kind:
			err.Kind = x
		case Op:
			err.Operation = x
		case Obj:
			err.Object = x
		case Code:
			err.Code = x
		case string:
			err.Details = err.Details + "; " + x
		case error:
			err.Err = x
		default:
			err.Details = err.Details + fmt.Sprintf("%v; ", x)
		}
	}
	err.Details = strings.Trim(err.Details, "; ")
	return err
}

// Error returns a string representation fo this Error.
func (err Error) Error() (ret string) {
	defer func() {
		ret = strings.Trim(ret, "\n- ;/")
	}()
	var e error = err
	for e != nil {
		if ret != "" {
			ret = ret + "\n"
		}
		ret = ret + " - "
		err2, ok := e.(Error)
		if !ok {
			ret = ret + e.Error()
			return
		}
		s := string(err2.Module)
		if err2.Function != "" {
			s = s + "/" + string(err2.Function)
		}
		if err2.Object != "" {
			s = s + " [" + string(err2.Object) + "]"
		}
		s = s + "\n   "
		if err2.Operation != "" {
			s = s + string(err2.Operation) + ": "
		}
		if err.Kind != "" {
			s = s + string(err2.Kind)
			if err.Code == "" {
				s = s + ":"
			}
			s = s + " "
		}
		if err2.Code != "" {
			s = s + "(" + string(err2.Code) + "): "
		}
		if err2.Details != "" {
			s = s + err2.Details
		}
		s = strings.Trim(s, " /\n:;")
		ret = ret + s
		e = err2.Err
	}
	return
}

// Unwrap returns the error this Error is wrapping. If no error is
// wrapped by this error nil will be returned.
func (err Error) Unwrap() error {
	return err.Err
}

const (

	// Some generic error kinds:
	NotFound        Kind = "not found"
	NotAllowed      Kind = "not allowed"
	IllegalArgument Kind = "illegal argument"
	IllegalValue    Kind = "illegal value"
	CreateErr       Kind = "can't create"
	ReadErr         Kind = "cant read"
	WriteErr        Kind = "can't write"
	DeleteErr       Kind = "can't delete"
	AlreadyExists   Kind = "already exists"
	ParseErr        Kind = "parsing failed"
	Failed          Kind = "operation failed"
	RecoveredPanic  Kind = "panic recovered"
)
