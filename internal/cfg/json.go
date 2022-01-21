package cfg

type jsonHolder struct {
	Command map[string]string      `json:"command"`
	Data    map[string]interface{} `json:"data"`
}
