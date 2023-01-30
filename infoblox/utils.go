package infoblox

import "fmt"

type ErrorCatcher struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func CheckErrorAndReact(e Exception) {
	if e != nil {
		panic(e)
	}
}

func (ec ErrorCatcher) Run() {
	if ec.Finally != nil {
		defer ec.Finally()
	}
	if ec.Catch != nil {
		defer func() {
			if e := recover(); e != nil {
				ec.Catch(e)
			}
		}()
	}
	if ec.Try == nil {
		panic("ErrorCatcher's 'Try' function must be defined")
	}

	ec.Try()
}

func safePtrValue(ptr interface{}) string {
	if ptr == nil {
		return fmt.Sprintf("nil")
	}

	switch t := ptr.(type) {
	case *int:
		if t == nil {
			return fmt.Sprintf("nil")
		}
		return fmt.Sprintf("%d", *t)
	case *uint32:
		if t == nil {
			return fmt.Sprintf("nil")
		}
		return fmt.Sprintf("%d", *t)
	case *bool:
		if t == nil {
			return fmt.Sprintf("nil")
		}
		return fmt.Sprintf("%t", *t)
	case *string:
		if t == nil {
			return fmt.Sprintf("nil")
		}
		return fmt.Sprintf("%s", *t)
	}

	panic("unsupported type")
}
