package fingerprint

const (
	GinPrint  = "gin"
	TRPCPrint = "tRPC"
	GoSwagger = "goswagger"
)

// 传入参数的key
const (
	ProductPathKey = "productPath"
)

// 满足其一
var (
	GinPrintPackages = []string{
		"github.com/gin-gonic/gin",
	}
	TRPCPrintPackage = []string{
		"git.code.oa.com/trpc-go/trpc-go",
		"trpc.group/trpc-go/trpc-go",
	}
	GoSwaggerPackage = []string{
		"github.com/go-openapi/swag",
	}
	AllTasks = map[string][]string{
		GinPrint:  GinPrintPackages,
		TRPCPrint: TRPCPrintPackage,
		GoSwagger: GoSwaggerPackage,
	}
)
