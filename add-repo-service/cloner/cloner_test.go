package cloner_test

import (
	"sync"
	"testing"

	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/cloner"
)

const (
	cloneLink = "https://github.com/roman-vrm/qatestlab.git"
	clonePath = "cloner-test-repo"
	branch    = "master"
)

var waitgroup sync.WaitGroup
var tt *testing.T

func TestCloneSuccess(t *testing.T) {
	tt = t
	waitgroup.Add(1)
	op := cloner.Options{
		Link:     cloneLink,
		Branch:   branch,
		Path:     clonePath,
		Callback: callbackSuccess,
	}
	cloner.StartClone(op)
	waitgroup.Wait()
	tt = nil
}

func TestCloneFailed(t *testing.T) {
	tt = t
	waitgroup.Add(1)
	op := cloner.Options{
		Link:     "cloneLink",
		Branch:   branch,
		Path:     clonePath,
		Callback: callbackFailed,
	}
	cloner.StartClone(op)
	waitgroup.Wait()
	tt = nil
}

func callbackSuccess(err error, op cloner.Options) {
	if err != nil {
		waitgroup.Done()
		tt.Fatal(err)
	}
	waitgroup.Done()
}

func callbackFailed(err error, op cloner.Options) {
	if err == nil {
		waitgroup.Done()
		tt.Fatal("clone without error")
	}
	waitgroup.Done()
}
