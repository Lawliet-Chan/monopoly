package web3

import (
	"errors"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"monopoly/common/logic"
	"monopoly/common/types"
	"net/http"
)

var GameNotFound = errors.New("game not found")

type GameManager struct {
	*tripod.Tripod
	games map[string]*logic.Game
}

func NewGameManager() *GameManager {
	gm := &GameManager{
		Tripod: tripod.NewTripod(),
		games:  make(map[string]*logic.Game),
	}
	gm.SetWritings(gm.CreateGame, gm.RollDice, gm.BuyProperty, gm.NextTurn)
	gm.SetReadings(gm.GetGameState)
	return gm
}

func (gm *GameManager) CreateGame(ctx *context.WriteContext) error {
	var req types.CreateGameRequest
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}
	gameID := logic.GenerateGameID() // 实现一个生成唯一ID的函数
	game := logic.NewGame(req.Players)

	gm.games[gameID] = game
	return nil
}

func (gm *GameManager) GetGameState(ctx *context.ReadContext) {
	gameID := ctx.GetString("gameId")
	game, exists := gm.games[gameID]
	if !exists {
		ctx.Json(http.StatusNotFound, context.H{"error": GameNotFound})
		return
	}
	ctx.JsonOk(types.GameStateResponse{
		CurrentPlayer: game.Current,
		Players:       game.Players,
		Board:         game.Board[:],
	})
}

func (gm *GameManager) RollDice(ctx *context.WriteContext) error {
	gameID := ctx.GetString("gameId")
	game, exists := gm.games[gameID]
	if !exists {
		return GameNotFound
	}
	err := game.MoveCurrentPlayer()
	if err != nil {
		return err
	}

	err = game.PayRent()
	if err != nil {
		return err
	}
	return ctx.EmitJsonEvent(types.ActionResponse{
		Success: true,
		Message: "购买成功",
	})
}

func (gm *GameManager) BuyProperty(ctx *context.WriteContext) error {
	gameID := ctx.GetString("gameId")
	game, exists := gm.games[gameID]
	if !exists {
		return GameNotFound
	}
	err := game.BuyProperty()
	if err != nil {
		return err
	}
	return ctx.EmitJsonEvent(types.ActionResponse{
		Success: true,
		Message: "回合结束",
	})
}

func (gm *GameManager) NextTurn(ctx *context.WriteContext) error {
	gameID := ctx.GetString("gameId")
	game, exists := gm.games[gameID]
	if !exists {
		return GameNotFound
	}
	game.NextPlayer()
	return ctx.EmitJsonEvent(types.ActionResponse{
		Success: true,
		Message: "回合结束",
	})
}
