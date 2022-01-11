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
	onGetNextKeyListener func() typ.Key
	onTerminateListener  func()
	onSaveConfigListener func(json typ.ConfigJSON)
	onLoadConfigListener func(typ.ConfigName) typ.ConfigJSON
	isOpen               bool
	configName           typ.ConfigName
	configNames          []typ.ConfigName
	jsonSchema           typ.ConfigJSONSchema
}

func (m *manager) SetOnGetNextKeyListener(f func() typ.Key) {
	m.onGetNextKeyListener = f
}

func (m *manager) SetOnTerminateListener(f func()) {
	m.onTerminateListener = f
}

func (m *manager) SetOnSaveConfigListener(f func(typ.ConfigJSON)) {
	m.onSaveConfigListener = f
}

func (m *manager) SetOnLoadConfigListener(f func(typ.ConfigName) typ.ConfigJSON) {
	m.onLoadConfigListener = f
}

func (m *manager) SetConfigNames(names []typ.ConfigName) {
	m.configNames = names
}

func (m *manager) ShowAlert(s string) {
	showAlert(s)
}

func (m *manager) ShowError(s string) {
	showError(s)
}

func (m *manager) SetJSONSchema(schema typ.ConfigJSONSchema) {
	m.jsonSchema = schema
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
		err := lorcaUI.Bind("__"+name+"__", f)
		if err != nil {
			panic(err)
		}
	}

	f(
		"ping",
		func() int {
			lg.Printf("Ping")
			return 1
		},
	)

	f(
		"getVersion",
		func() string {
			fmt.Println("Call GetVersion")
			r := cnst.VERSION
			fmt.Println("Return GetVersion")
			return r
		},
	)

	f(
		"getSchema",
		func() string {
			fmt.Println("Call GetScehm")
			r := string(m.jsonSchema)
			fmt.Println("Return GetScDS")
			return r
		},
	)

	f(
		"openLink",
		func(url string) {
			_ = exec.Command(
				"rundll32", "url.dll,FileProtocolHandler", "https://github.com/wirekang/mouseable",
			).Start()
		},
	)

	f(
		"terminate",
		func() {
			m.onTerminateListener()
		},
	)

	f(
		"getConfigNames",
		func() []string {
			r := make([]string, len(m.configNames))
			for i := range m.configNames {
				r[i] = string(m.configNames[i])
			}
			return r
		},
	)

	f(
		"getConfig",
		func(name string) string {
			return string(m.onLoadConfigListener(typ.ConfigName(name)))
		},
	)

	f(
		"saveConfig",
		func(json string) {
			m.onSaveConfigListener(typ.ConfigJSON(json))
		},
	)
}
