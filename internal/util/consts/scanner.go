package consts

const (
	ArgsSize   = 2
	BeginLevel = 0
)

// function names
const (
	ContainFunc    = "contain"
	BeginStrFunc   = "beginStr"
	EndStrFunc     = "endStr"
	RegFunc        = "reg"
	BeginWithLower = "beginWithLower"
	BeginWithUpper = "beginWithUpper"
	EqualFunc      = "equal"
)

// mode
const (
	StrMode string = "str"
	GoMode  string = "go"
	AiMode  string = "ai"
)

// 用于传参数
const (
	IsUseCtxKey     = "isUseCtx"
	AiConfigKey     = "aiConfig"
	IsReturnBoolKey = "isReturnBool"
)

// 用于指定file path
const (
	GinRule       = "rule/find_gin_api.yaml"
	GoSwaggerRule = "rule/find_go_swagger_api.yaml"
	TrpcRule      = "rule/find_trpc_api.yaml"
)
