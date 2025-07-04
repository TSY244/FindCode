package scanner

import (
	"ScanIDOR/internal/pkg/fcFlag"
	"sync"
)

type Env struct {
	Result          map[string][]string
	funcCache       Cache
	ApiCache        Cache
	ApiCacheMap     map[string]*cacheUnit
	JudgedCache     map[string]bool
	CodeCache       map[string][]*cacheUnit
	nameJudgedCache map[string]bool
	modeJudgeCache  map[string]struct{}
	AiBoolResult    map[string]AiBoolResultUnit
	FuncCacheMap    map[string]*cacheUnit
	AiCycle         int
}

func NewEnv() *Env {
	env := &Env{
		Result:          make(map[string][]string),
		funcCache:       make(Cache),
		FuncCacheMap:    make(map[string]*cacheUnit),
		ApiCache:        make(Cache),
		ApiCacheMap:     make(map[string]*cacheUnit),
		JudgedCache:     make(map[string]bool),
		CodeCache:       make(map[string][]*cacheUnit),
		nameJudgedCache: make(map[string]bool),
		modeJudgeCache:  make(map[string]struct{}),
		AiBoolResult:    make(map[string]AiBoolResultUnit),
	}
	env.AiCycle = fcFlag.AiCycle
	return env
}

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

	// CodeCache 代码的缓存
	CodeCache = make(map[string][]*cacheUnit)

	// nameJudgedCache 有鉴权框架返回存放false
	nameJudgedCache = make(map[string]bool)

	// modeCache 存放扫描模式
	modeCache = make(map[string]struct{})

	// AiBoolResult 存放ai json 格式的扫描结果
	AiBoolResult = make(map[string]AiBoolResultUnit) // path -> funcName -> result

	// 为了后端项目兼容，使用读写锁
	locker = new(sync.RWMutex)
)
