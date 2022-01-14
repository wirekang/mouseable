package ui

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"

	"github.com/zserge/lorca"

	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/typ"
)

func New() typ.UIManager {
	mustChrome()
	return &manager{}
}

type manager struct {
	onGetNextKeyListener            func() typ.Key
	onTerminateListener             func()
	onSaveConfigListener            func(typ.ConfigName, typ.ConfigJSON) error
	onLoadConfigListener            func(typ.ConfigName) (typ.ConfigJSON, error)
	onLoadConfigSchemaListener      func() typ.ConfigJSONSchema
	onLoadConfigNamesListener       func() ([]typ.ConfigName, error)
	onLoadAppliedConfigNameListener func() (typ.ConfigName, error)
	onApplyConfigListener           func(typ.ConfigName) error
	isOpen                          bool
}

func (m *manager) SetOnApplyConfigNameListener(f func(name typ.ConfigName) error) {
	m.onApplyConfigListener = f
}

func (m *manager) SetOnLoadAppliedConfigNameListener(f func() (typ.ConfigName, error)) {
	m.onLoadAppliedConfigNameListener = f
}

func (m *manager) SetOnLoadConfigSchemaListener(f func() typ.ConfigJSONSchema) {
	m.onLoadConfigSchemaListener = f
}

func (m *manager) SetOnLoadConfigNamesListener(f func() ([]typ.ConfigName, error)) {
	m.onLoadConfigNamesListener = f
}

func (m *manager) SetOnGetNextKeyListener(f func() typ.Key) {
	m.onGetNextKeyListener = f
}

func (m *manager) SetOnTerminateListener(f func()) {
	m.onTerminateListener = f
}

func (m *manager) SetOnSaveConfigListener(f func(typ.ConfigName, typ.ConfigJSON) error) {
	m.onSaveConfigListener = f
}

func (m *manager) SetOnLoadConfigListener(f func(typ.ConfigName) (typ.ConfigJSON, error)) {
	m.onLoadConfigListener = f
}

func (m *manager) ShowAlert(s string) {
	showAlert(s)
}

func (m *manager) ShowError(s string) {
	showError(s)
}

func (m *manager) Open() {
	m.openUI()
}

func mustChrome() {
	if lorca.LocateChrome() == "" {
		panic("Chromium browser not found. Mouseable can't render GUI. Please install Chrome or Edge.")
	}
}

func (m *manager) openUI() {
	if m.isOpen {
		lg.Printf("Window is already open.")
		return
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	go func() {
		err = http.Serve(listener, http.FileServer(http.FS(cnst.FrontFS)))
		if err != nil {
			lg.Errorf("http.Serve: %v", err)
		}
	}()

	host := "http://" + listener.Addr().String()
	lg.Printf("Host: %s", host)
	lorcaUI, err := lorca.New(host, "", 800, 800, "--disable-features=Translate")
	if err != nil {
		panic(err)
	}

	m.bindLorca(lorcaUI)
	m.isOpen = true
	defer func() {
		m.isOpen = false
		err = lorcaUI.Close()
		if err != nil {
			lg.Errorf("ui.Close(): %v", err)
		}

		err = listener.Close()
		if err != nil {
			lg.Errorf("listener.Close(): %v", err)
		}
		lg.Printf("Close UI")
	}()

	isTerminate := false
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	select {
	case <-sigChan:
		isTerminate = true
	case <-lorcaUI.Done():
	}

	if isTerminate {
		m.onTerminateListener()
	}
}

func (m *manager) bindLorca(lorcaUI lorca.UI) {
	f := func(name string, f interface{}) {
		fmt.Println("bind", name)
		err := lorcaUI.Bind("__"+name, f)
		if err != nil {
			panic(err)
		}
	}
	f("terminate", m.onTerminateListener)
	f("loadConfigNames", m.onLoadConfigNamesListener)
	f("loadConfig", m.onLoadConfigListener)
	f("saveConfig", m.onSaveConfigListener)
	f("loadSchema", m.onLoadConfigSchemaListener)
	f("getNextKey", m.onGetNextKeyListener)
	f("loadAppliedConfigName", m.onLoadAppliedConfigNameListener)
	f("applyConfig", m.onApplyConfigListener)

	f("ping", func() int { return 1 })
	f("getVersion", func() string { return cnst.VERSION })
	f(
		"openLink", func(url string) {
			_ = exec.Command(
				"rundll32", "url.dll,FileProtocolHandler", "https://github.com/wirekang/mouseable",
			).Start()
		},
	)
}
