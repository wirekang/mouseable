package cfg

type jsonHolder struct {
	Command map[string]string      `json:"Command"`
	Data    map[string]interface{} `json:"Data"`
}
