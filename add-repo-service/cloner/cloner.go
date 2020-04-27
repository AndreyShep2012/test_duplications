package cloner

import (
	gitlib "gitlab.qarea.org/jiraquality/cq-web-backend/git-lib"
)

//Callback is called when clone is completed
type Callback func(error, Options)

//Options -
type Options struct {
	Link     string
	Path     string
	Branch   string
	Folder   string
	RepoID   string
	SSHKey   []byte
	Callback Callback
}

//StartClone -
func StartClone(o Options) {
	go processClone(o)
}

func processClone(o Options) {
	o.Callback(gitlib.CloneRepo(o.Link, o.Branch, o.Path, o.SSHKey), o)
}
