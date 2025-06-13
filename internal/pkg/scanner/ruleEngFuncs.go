package scanner

import (
	"fmt"
	"regexp"
	"strings"
)

func base(f func(args ...interface{}) (interface{}, error),
	args ...interface{}) (interface{}, error) {
	if ArgsSize != len(args) {
		return false, ArgsSizeNotEqualErr
	}
	return f(args...)
}

func Contain(args ...interface{}) (interface{}, error) {
	return base(func(args ...interface{}) (interface{}, error) {
		return strings.Contains(fmt.Sprint(args[1]), fmt.Sprint(args[0])), nil
	}, args...)
}

func EndStr(args ...interface{}) (interface{}, error) {
	return base(func(args ...interface{}) (interface{}, error) {
		return strings.HasSuffix(fmt.Sprint(args[1]), fmt.Sprint(args[0])), nil
	}, args...)
}

func BeginStr(args ...interface{}) (interface{}, error) {
	return base(func(args ...interface{}) (interface{}, error) {
		return strings.HasPrefix(fmt.Sprint(args[1]), fmt.Sprint(args[0])), nil
	}, args...)
}

func Reg(args ...interface{}) (interface{}, error) {
	return base(func(args ...interface{}) (interface{}, error) {
		re := regexp.MustCompile(fmt.Sprint(args[0]))
		return re.MatchString(fmt.Sprint(args[1])), nil
	}, args...)
}

func BeginWithLowerCase(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, FuncParamsNotEqualErr
	}
	data := args[0].(string)
	return data[0] == strings.ToLower(data)[0], nil
}

func BeginWithUpperCase(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, FuncParamsNotEqualErr
	}
	data := args[0].(string)
	return data[0] == strings.ToUpper(data)[0], nil
}

func Equal(args ...interface{}) (interface{}, error) {
	return base(func(args ...interface{}) (interface{}, error) {
		return args[0] == args[1], nil
	}, args...)
}
