package util

import (
	"context"
	"testing"
)

func TestCloneRepository(t *testing.T) {
	ctx := context.Background()
	if repo, err := CloneRepository(ctx, "https://github.com/TSY244/FindCode.git"); err != nil {
		t.Error(err)
	} else {
		t.Log("clone success:", repo)
	}
}
