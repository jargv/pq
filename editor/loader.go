package editor

import (
	"io/ioutil"
	_ "log"
	"os"
	"os/user"
	"path"
	"prio/prio-cli/project"
	"strings"
)

type loader struct {
	Project  *project.Project
	filename string
	isDir    bool
}

func (e *loader) OpenTopLevel() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	dir := path.Join(usr.HomeDir, ".prio")
	err = e.Open(dir)
	//if it's a path err, make dir and try again
	if pe, ok := err.(*os.PathError); ok {
		if pe.Path == dir {
			if err = os.Mkdir(dir, 0700); err != nil {
				return err
			}
			err = e.Open(dir)
		}
	}

	for err == nil && e.isDir {
		if e.Project.NTasks() == 0 {
			break
		}

		task := e.Project.Get(0)
		dir := path.Dir(e.filename)
		path := path.Join(dir, task)
		err = e.Open(path)
	}

	return err
}

func (e *loader) Open(filename string) error {
	info, err := os.Stat(filename)
	//path error is fine here, it just means project is new
	if _, ok := err.(*os.PathError); err != nil && !ok {
		return err
	}

	if info != nil && info.IsDir() {
		e.isDir = true
		indexname := path.Join(filename, ".index")
		return e.openFile(indexname)
	}

	e.isDir = false
	return e.openFile(filename)
}

func (e *loader) openFile(filename string) error {
	if e.Project != nil {
		e.Save()
		e.Project = nil
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			contents = []byte{}
		} else {
			return err
		}
	}

	items := strings.Split(string(contents), "\n")
	e.filename = filename
	e.Project = project.New(items)

	return nil
}

func (e *loader) Save() error {
	if e.Project == nil {
		return nil
	}
	items := e.Project.Items()
	return ioutil.WriteFile(e.filename, []byte(strings.Join(items, "\n")), 0777)
}
