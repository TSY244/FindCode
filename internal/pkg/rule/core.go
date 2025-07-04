package rule

import (
	"ScanIDOR/internal/config"
	"ScanIDOR/internal/pkg/fcFlag"
	"ScanIDOR/internal/pkg/global"
	"ScanIDOR/internal/util/consts"
	"ScanIDOR/pkg/logger"
	"ScanIDOR/utils/util"
)

func LoadRuleWithFrame(frames []string) []Rule {
	var rules []Rule
	for _, frame := range frames {
		configPath := global.RuleMap[frame]
		var r Rule
		if err := util.LoadYaml(configPath, &r); err != nil {
			logger.Fatal(err)
		}
		rules = append(rules, r)
	}
	return rules
}

func LoadRule(configPath string, c *config.Config) Rule {
	var r Rule
	if err := util.LoadYaml(configPath, &r); err != nil {
		logger.Fatal(err)
	}
	loadEnv(&r, c)
	return r
}

func loadEnv(r *Rule, c *config.Config) {
	if fcFlag.GoTarget != "" {
		r.GoModeTargetRule.Rule = fcFlag.GoTarget
	}
	if fcFlag.AiMode == true {
		r.Mode = append(r.Mode, consts.AiMode)
		r.AiConfig = c.AiConfig
	}
}
