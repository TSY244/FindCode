package scanner

import "ScanIDOR/pkg/ruleEng"

func matchStr(rule, target string) (bool, error) {
	return ruleEng.Run(rule, target)
}

func matchStrSplice(targets, rules []string) (bool, error) {
	for i, rule := range rules {
		if len(targets) <= i {
			return true, nil
		}
		if ret, err := matchStr(rule, targets[i]); err != nil {
			return false, err
		} else if !ret {
			return false, nil
		}
	}
	return true, nil
}
