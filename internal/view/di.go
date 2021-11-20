package view

var DI struct {
	LoadData func() (map[string][]uint32, map[string]string, error)
	SaveData func(map[string][]uint32, map[string]string) error
}
