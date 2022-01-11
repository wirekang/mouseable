package io

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/juju/fslock"
	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/typ"
)

func New() typ.IOManager {
	rd, cd := initDirs()
	return &manager{
		rootDir:   rd,
		configDir: cd,
	}
}

type manager struct {
	rootDir   string
	configDir string
	lock      *fslock.Lock
}

func (i *manager) LoadConfigNames() (rst []typ.ConfigName) {
	rst = make([]typ.ConfigName, 0)
	infos, err := ioutil.ReadDir(i.configDir)
	if err != nil {
		return
	}

	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		rst = append(rst, typ.ConfigName(info.Name()))
	}
	return
}

func (i *manager) SaveConfig(name typ.ConfigName, data typ.ConfigJSON) (err error) {
	err = ioutil.WriteFile(filepath.Join(i.configDir, string(name)), []byte(data), fs.ModePerm)
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (i *manager) LoadConfig(name typ.ConfigName) (data typ.ConfigJSON, err error) {
	bytes, err := ioutil.ReadFile(filepath.Join(i.configDir, string(name)))
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	data = typ.ConfigJSON(bytes)
	return
}

func (i *manager) Lock() {
	lockFile := i.rootDir + "\\lockfile"
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

func initDirs() (rootDir string, configDir string) {
	if runtime.GOOS == "windows" {
		rootDir = filepath.Join(os.Getenv("APPDATA"), "mouseable")
		if cnst.IsDev {
			rootDir += "_dev"
		}
		configDir = filepath.Join(rootDir, "configs")

		if isNotExists(configDir) {
			_ = os.MkdirAll(configDir, os.ModeDir)
			initDefaultConfigs(configDir)
		}
		return
	}
	panic(fmt.Sprintf("%s not supported.", runtime.GOOS))
}

func isNotExists(dir string) bool {
	_, err := os.Stat(dir)
	return os.IsNotExist(err)
}

func initDefaultConfigs(configDir string) {
	lg.Printf("Init default configs")

	entries, err := fs.ReadDir(cnst.DefaultConfigsFS, ".")
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			lg.Errorf("%s is directory", entry.Name())
			continue
		}

		src, err := cnst.DefaultConfigsFS.Open(entry.Name())
		if err != nil {
			lg.Errorf("DefaultConfigsFS.Open: %v", err)
			continue
		}

		dst, err := os.Create(filepath.Join(configDir, entry.Name()))
		if err != nil {
			lg.Errorf("os.Open: %v", err)
			continue
		}
		_, err = io.Copy(dst, src)
		if err != nil {
			lg.Errorf("dst.Write: %v", err)
			continue
		}

		src.Close()
		dst.Close()
	}
}
