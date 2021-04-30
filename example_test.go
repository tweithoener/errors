package errors_test

import (
	"fmt"

	"github.com/tweithoener/errors"
)

func ExampleT() {
	type handler func(b string) error
	makeHandler := func(a string) handler {

		// Create an error template that is good for all handlers creater by this function
		etpl := errors.T(
			errors.Mod("example"),
			errors.Func("handler"),
			errors.Obj(a))
		return func(b string) error {
			if b == "" {
				// crete a specific error from template and additional attributes
				return etpl.E(errors.Op("argument check"), errors.Kind("b is empty string"), errors.Code("E123"))
			}
			if a != b {
				// create a different error from the same template
				return etpl.E(errors.Op("compare"), errors.Kind("not equal"), errors.Code("E456"), a, b)
			}
			return nil
		}
	}

	handleTest := makeHandler("test")
	handleCheck := makeHandler("check")

	if err := handleTest(""); err != nil {
		fmt.Println(err.Error())
	}
	if err := handleTest("test"); err != nil {
		fmt.Println(err.Error())
	}
	if err := handleTest("not test"); err != nil {
		fmt.Println(err.Error())
	}
	if err := handleCheck("check"); err != nil {
		fmt.Println(err.Error())
	}
	if err := handleCheck("test"); err != nil {
		fmt.Println(err.Error())
	}
	// Output:
	// example/handler [test]
	//    argument check: b is empty string (E123)
	// example/handler [test]
	//    compare: not equal (E456): test; not test
	// example/handler [check]
	//    compare: not equal (E456): check; test
}

func ExampleE() {
	// a very basic error. it will only consist of the provided message.
	err1 := errors.E("file does not exist")
	fmt.Println("Error 1:\n" + err1.Error())

	// another simple error wrapping the above error
	err2 := errors.E("can't read config", err1)
	fmt.Println("Error 2:\n" + err2.Error())

	// and one more ... with a lot more more attributes. See how the errors
	// are printed like a stacktrace.
	err3 := errors.E(errors.Mod("example"), errors.Func("startup"), errors.Op("configure"), errors.Kind("configure failed"), err2)
	fmt.Println("Error 3:\n" + err3.Error())

	// you can wrap arbitrary errors
	otherErr := fmt.Errorf("other error")
	err4 := errors.E(otherErr, errors.Kind("wrapper"))
	fmt.Println("Error 4:\n" + err4.Error())

	// arbitrary data is allowed as error detail. it will be converted into a string using fmt.Sprintf("%v", ...)
	a := struct {
		A int
		B int
	}{1, 2}
	err5 := errors.E(errors.Kind("illegal value"), a)
	fmt.Println("Error 5:\n" + err5.Error())

	// Output:
	// Error 1:
	// file does not exist
	// Error 2:
	// can't read config
	//  - file does not exist
	// Error 3:
	// example/startup
	//    configure: configure failed
	//  - can't read config
	//  - file does not exist
	// Error 4:
	// wrapper
	//  - other error
	// Error 5:
	// illegal value: {1 2}
}

func ExampleTemplate_E() {
	// create a template
	etpl := errors.T(errors.Func("Template.T example"))

	// and two errors based on this template. They will share the same Func attribute, defined in the template.
	err1 := etpl.E(errors.Op("operation 1"))
	err2 := etpl.E(errors.Op("operation 2"))

	fmt.Printf("%t %s", err1.Function == err2.Function, err2.Function)

	// Also see the example for errors.T() for a more detailed example!

	// Output:
	// true Template.T example
}

func ExampleTemplate_T() {
	// create a template
	etpl1 := errors.T(errors.Func("Template.T example"))

	// and base another template on this one
	etpl2 := etpl1.T(errors.Obj("etpl2"))

	// this errors will have the Func attribute (from the first template) and the Op attributes set.
	err1 := etpl1.E(errors.Op("operation 1"))
	// additionally this error will have the Obj attribute set (from the second template)
	err2 := etpl2.E(errors.Op("operation 2"))

	fmt.Printf("%t %s\n", err1.Function == err2.Function, err2.Function)
	fmt.Printf("%t %s\n", err1.Object == err2.Object, err2.Object)

	// Also see the example for errors.T() for a more detailed example!

	// Output:
	// true Template.T example
	// false etpl2
}
