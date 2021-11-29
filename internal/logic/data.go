package logic

import (
	"github.com/wirekang/mouseable/internal/def"
)

type dataCache struct {
	int   int
	float float64
	bool  bool
	string
}

func SetConfig(config def.Config) {
	mutex.Lock()
	cacheData(config.DataMap)
	functionMap = config.FunctionMap
	mutex.Unlock()
}

func cacheData(dm def.DataMap) {
	cachedDataMap = make(map[*def.DataDefinition]dataCache, len(dm))
	for dd, v := range dm {
		cachedDataMap[dd] = dataCache{
			int:    v.Int(),
			float:  v.Float(),
			bool:   v.Bool(),
			string: string(v),
		}
	}

}
