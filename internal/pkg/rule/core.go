package rule

import (
	"ScanIDOR/internal/pkg/global"
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

func LoadRule(configPath string) Rule {
	var r Rule
	if err := util.LoadYaml(configPath, &r); err != nil {
		logger.Fatal(err)
	}
	return r
}
