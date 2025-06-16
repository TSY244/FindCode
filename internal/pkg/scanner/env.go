package scanner

var (
	// Result 用于存放并发时候的数据
	Result = make(map[string][]string)
	// funcCache 用于存放非api 的函数的地方
	funcCache = make(Cache)

	// funcCacheMap 用于存放普通代码
	funcCacheMap = make(map[string]*cacheUnit)

	// cacheUnit 用于存放api 的地方
	apiCache = make(Cache)

	//JudgedCache 判断过的缓存，有鉴权框架返回存放false
	JudgedCache = make(map[string]bool)

	// CodeCache
	CodeCache = make(map[string][]*cacheUnit)

	// nameJudgedCache 有鉴权框架返回存放false
	nameJudgedCache = make(map[string]bool)

	// modeCache 存放扫描模式
	modeCache = make(map[string]struct{})

	// AiResult 存放ai 扫描结果
	AiResult = make(map[string]AiResultUnit) // path -> funcName -> result
)
