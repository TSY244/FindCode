package rule

type Rule struct {
	TaskName          string         `mapstructure:"task_name"`
	GoModeTargetRule  TargetRule     `mapstructure:"go_mode_target_rule"`
	StrModeTargetRule TargetRule     `mapstructure:"str_mode_target_rule"`
	Path              Path           `mapstructure:"path"`
	File              File           `mapstructure:"file"`
	FuncRules         []FuncRuleUnit `mapstructure:"func_rule"`
	Mode              []string       `mapstructure:"mode"`
}

type FuncRuleUnit struct {
	FuncNameRule   *FuncNameRule  `mapstructure:"func_name"`
	ParamNameRule  *ParamNameRule `mapstructure:"param_name"`
	ParamTypeRule  *ParamTypeRule `mapstructure:"param_type"`
	ReturnTypeRule *ReturnType    `mapstructure:"return_type"`
	RecvTypeRule   *RecvTypeRule  `mapstructure:"recv_type"`
	RecvNameRule   *RecvNameRule  `mapstructure:"recv_name"`
}

type RecvNameRule struct {
	Rule string `mapstructure:"rule"`
}

type RecvTypeRule struct {
	Rule string `mapstructure:"rule"`
}

type ReturnType struct {
	Rules []ReturnTypeRule `mapstructure:"rules"`
}

type ReturnTypeRule struct {
	Size  int      `mapstructure:"size"`
	Rules []string `mapstructure:"rules"`
}

type TargetRule struct {
	Rule string `mapstructure:"rule"`
}

type FuncNameRule struct {
	Rule string `mapstructure:"rule"`
}

type ParamNameRule struct {
	Rule []ParamUnit `mapstructure:"rules"`
}

type ParamTypeRule struct {
	Rule []ParamUnit `mapstructure:"rules"`
}

type Path struct {
	Rule     string `mapstructure:"rule"`
	DeepSize int    `mapstructure:"deepsize"`
}

type File struct {
	Rule string `mapstructure:"rule"`
}

type ParamUnit struct {
	Size  int      `mapstructure:"size"`
	Rules []string `mapstructure:"rules"`
}
