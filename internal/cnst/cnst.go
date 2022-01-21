package cnst

import (
	"io/fs"
)

var VERSION string
var IsDev bool
var AssetFS fs.FS
var FrontFS fs.FS
var DefaultConfigsFS fs.FS
