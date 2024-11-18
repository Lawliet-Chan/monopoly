package web2

import (
	"github.com/gin-gonic/gin"
	"monopoly/common/logic"
	"monopoly/common/types"
	"net/http"
	"sync"
)

// GameManager 管理多个游戏实例
type GameManager struct {
	games map[string]*logic.Game // gameID -> Game
	mu    sync.RWMutex
}

// 创建新的游戏管理器
func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[string]*logic.Game),
	}
}

// Handler 处理HTTP请求的结构体
type Handler struct {
	manager *GameManager
}

// NewHandler 创建新的Handler
func NewHandler(manager *GameManager) *Handler {
	return &Handler{manager: manager}
}

// SetupRoutes 设置路由
func (h *Handler) SetupRoutes(router *gin.Engine) {
	// 游戏管理相关接口
	router.POST("/games", h.CreateGame)
	router.GET("/games/:gameId", h.GetGameState)

	// 游戏操作相关接口
	game := router.Group("/games/:gameId")
	{
		game.POST("/roll", h.RollDice)
		game.POST("/buy", h.BuyProperty)
		game.POST("/next", h.NextTurn)
	}
}

// CreateGame 创建新游戏
func (h *Handler) CreateGame(c *gin.Context) {
	var req types.CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gameID := logic.GenerateGameID() // 实现一个生成唯一ID的函数
	game := logic.NewGame(req.Players)

	h.manager.mu.Lock()
	h.manager.games[gameID] = game
	h.manager.mu.Unlock()

	c.JSON(http.StatusCreated, types.CreateGameResponse{GameID: gameID})
}

// GetGameState 获取游戏状态
func (h *Handler) GetGameState(c *gin.Context) {
	gameID := c.Param("gameId")

	h.manager.mu.RLock()
	game, exists := h.manager.games[gameID]
	h.manager.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	c.JSON(http.StatusOK, types.GameStateResponse{
		CurrentPlayer: game.Current,
		Players:       game.Players,
		Board:         game.Board,
	})
}

// RollDice 掷骰子并移动
func (h *Handler) RollDice(c *gin.Context) {
	gameID := c.Param("gameId")

	h.manager.mu.Lock()
	game, exists := h.manager.games[gameID]
	h.manager.mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	if err := game.MoveCurrentPlayer(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 自动触发租金支付
	if err := game.PayRent(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.ActionResponse{
		Success: true,
		Message: "移动完成",
	})
}

// BuyProperty 购买当前位置的地产
func (h *Handler) BuyProperty(c *gin.Context) {
	gameID := c.Param("gameId")

	h.manager.mu.Lock()
	game, exists := h.manager.games[gameID]
	h.manager.mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	if err := game.BuyProperty(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.ActionResponse{
		Success: true,
		Message: "购买成功",
	})
}

// NextTurn 结束当前回合
func (h *Handler) NextTurn(c *gin.Context) {
	gameID := c.Param("gameId")

	h.manager.mu.Lock()
	game, exists := h.manager.games[gameID]
	h.manager.mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	game.NextPlayer()

	c.JSON(http.StatusOK, types.ActionResponse{
		Success: true,
		Message: "回合结束",
	})
}
