package limit

import (
	"sync"
)

type LimiterCache struct {
	Data sync.Map
}

var IpCache *LimiterCache


func init() {

	IpCache = &LimiterCache{}

}



