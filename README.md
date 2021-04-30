# errors
Golang package to quickly create descriptive, wrapped errors. Features error templates.

See the [package documentation](https://pkg.go.dev/github.com/tweithoener/errors#section-documentation) for documentation and usage examples.

## Usage Example
A basic example might loog like this

```golang
    package yourpkg

    import "github.com/tweithoener/errors"

    var (
        etpl        = errors.T(errors.Mod("your package"))   // template for the whole pkg.
        etSomeFunc  = etpl.T(errors.Func("some function"))   // templates for specific functions
        etOtherFunc = etpl.T(errors.Func("other function"))  // with module and function information
    )

    func someFunc() error {
        err := doStuff()
        if err != nil {
            // error creation using the template from above
            // error now holds information on the module, the
            // function, the current operation and the cause.
            return etSomeFunc.E(error.Op("doing stuff"), err)
        }
        err = moreStuff()
        if err != nil {
            // module and function information as above (same template).
            // operation and cause specific to this error.
            return etSomeFunc.E(error.Op("doing more stuff"), err)
        }
        return nil
    }

    func otherFunc() error {
        // ...
        return etOtherFunc.E(...) 
    }
```

## Installation
Install using 
```
    go get github.com/tweithoener/errors
```

or require this package as a module in your go.mod file.

## License
MIT License. See LICENSE file.
