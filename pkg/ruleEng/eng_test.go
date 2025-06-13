package ruleEng

import (
	"testing"
)

func TestProcessRule(t *testing.T) {
	ret := processRule(`!contain("check")&&!contain("Check")`)
	t.Logf("ret:%v", ret)
}
