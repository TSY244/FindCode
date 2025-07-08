package scanner

import "errors"

var (
	FuncParamsNotEqualErr = errors.New("funcParams not equal")
	FuncParamsNotEqual    = errors.New("funcParams not equal")
)

var (
	ArgsSizeNotEqualErr = errors.New("args size not equal")
)

var (
	NoRecv  = errors.New("no receiver")
	RecvErr = errors.New("recv error")
)
