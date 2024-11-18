package types

import "monopoly/common/logic"

// 请求和响应的结构体
type CreateGameRequest struct {
	Players []string `json:"players" binding:"required,min=2"`
}

type CreateGameResponse struct {
	GameID string `json:"gameId"`
}

type GameStateResponse struct {
	CurrentPlayer int             `json:"currentPlayer"`
	Players       []*logic.Player `json:"players"`
	Board         []*logic.Place  `json:"board"`
}

type ActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
