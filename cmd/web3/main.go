package main

import (
	"github.com/yu-org/yu/apps/poa"
	"github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/startup"
	"monopoly/web3"
)

func main() {
	yuCfg := config.InitDefaultCfg()

	poaCfg := poa.DefaultCfg(0)
	poaTripod := poa.NewPoa(poaCfg)
	gameTripod := web3.NewGameManager()

	startup.InitDefaultKernel(yuCfg).
		WithTripods(
			poaTripod,
			gameTripod,
		).Startup()
}
