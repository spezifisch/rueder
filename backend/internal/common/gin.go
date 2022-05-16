package common

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

func GinSetTrustedProxies(engine *gin.Engine, trustedProxies []string) {
	// set trusted proxies, see https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	if len(trustedProxies) == 0 {
		// don't trust any client IP headers by default
		if err := engine.SetTrustedProxies(nil); err != nil {
			panic("failed setting trusted proxies to nil")
		}
		log.Info("not trusting any proxy")
	} else {
		if err := engine.SetTrustedProxies(trustedProxies); err != nil {
			panic("failed setting trusted proxies")
		}
		log.Infof("trusted proxies: %s", trustedProxies)
	}
}
