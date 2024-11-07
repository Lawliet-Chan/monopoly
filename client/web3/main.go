package main

import (
	"github.com/yu-org/go-yu-sdk/pkg"
	"monopoly/common/types"
)

func main() {
	cli := pkg.NewClient("http://localhost:7999").WithLeiPrice(1)
	err := cli.WriteChain(
		"gamemanager", "CreateGame",
		types.CreateGameRequest{
			Players: []string{"Alice"},
		})
	if err != nil {
		panic(err)
	}

}
