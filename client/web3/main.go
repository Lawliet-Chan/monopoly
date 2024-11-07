package main

import (
	"github.com/yu-org/go-yu-sdk/pkg"
	"github.com/yu-org/yu/core/keypair"
	"monopoly/common/types"
)

func main() {
	pubkey, privkey := keypair.GenEdKeyWithSecret([]byte("Alice"))
	cli := pkg.NewClient("http://localhost:7999").WithLeiPrice(1).WithKeys(privkey, pubkey)
	err := cli.WriteChain(
		"gamemanager", "CreateGame",
		types.CreateGameRequest{
			Players: []string{"Alice"},
		})
	if err != nil {
		panic(err)
	}

}
