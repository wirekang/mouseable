package io

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/juju/fslock"
	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/typ"
)

func New() typ.IOManager {
	dataDir := getDataDir()
	return &manager{
		dataDir: dataDir,
	}
}

type manager struct {
	dataDir string
	lock    *fslock.Lock
}

func (i *manager) LoadNames() (rst []string) {
	rst = make([]string, 0)
	infos, err := ioutil.ReadDir(i.dataDir)
	if err != nil {
		return
	}

	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		rst = append(rst, info.Name())
	}
	return
}

func (i *manager) Save(name typ.ConfigName, data typ.ConfigJSON) (err error) {
	err = ioutil.WriteFile(filepath.Join(i.dataDir, string(name)), []byte(data), fs.ModePerm)
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (i *manager) Load(name typ.ConfigName) (data typ.ConfigJSON, err error) {
	bytes, err := ioutil.ReadFile(filepath.Join(i.dataDir, string(name)))
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	data = typ.ConfigJSON(bytes)
	return
}

func (i *manager) Lock() {
	lockFile := i.dataDir + "\\lockfile"
	i.lock = fslock.New(lockFile)
	err := i.lock.TryLock()
	if err != nil {
		if errors.Is(err, fslock.ErrLocked) {
			panic("Mouseable is already running. Please check tray icon.")
		}

		panic(err)
	}

	return
}

func (i *manager) Unlock() {
	_ = i.lock.Unlock()
}

func getDataDir() (dataDir string) {
	if runtime.GOOS == "windows" {
		dataDir = filepath.Join(os.Getenv("APPDATA"), "mouseable")
		if cnst.IsDev {
			dataDir += "_dev"
		}
		return
	}
	panic(fmt.Sprintf("%s not supported.", runtime.GOOS))
}
