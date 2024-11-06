package logic

import (
	"errors"
	"math/rand"
)

// 地块类型
type Space struct {
	ID       int
	Name     string
	Price    int     // 购买价格
	Rent     int     // 租金
	Owner    *Player // 所有者
	Position int     // 在棋盘上的位置
}

// 玩家结构
type Player struct {
	ID       int
	Name     string
	Money    int
	Position int
	Property []*Space // 拥有的地产
}

// 游戏结构
type Game struct {
	Board   [40]*Space // 棋盘，40个格子
	Players []*Player
	Current int // 当前玩家索引
}

// 创建新游戏
func NewGame(playerNames []string) *Game {
	game := &Game{
		Players: make([]*Player, len(playerNames)),
	}

	// 初始化玩家
	for i, name := range playerNames {
		game.Players[i] = &Player{
			ID:       i,
			Name:     name,
			Money:    1500, // 初始资金
			Position: 0,
			Property: make([]*Space, 0),
		}
	}

	// 初始化棋盘
	for i := 0; i < 40; i++ {
		if i%5 != 0 { // 每隔5个空格设置一个可购买的地产
			game.Board[i] = &Space{
				ID:       i,
				Name:     "地产" + string(rune('A'+i)),
				Price:    (i + 1) * 100,
				Rent:     (i + 1) * 50,
				Position: i,
			}
		}
	}

	return game
}

// 掷骰子
// TODO: 将改成链上VRF随机数实现
func (g *Game) RollDice() int {
	return rand.Intn(6) + 1 + rand.Intn(6) + 1
}

// 移动当前玩家
func (g *Game) MoveCurrentPlayer() error {
	if len(g.Players) == 0 {
		return errors.New("没有玩家")
	}

	player := g.Players[g.Current]
	steps := g.RollDice()

	// 更新位置，如果超过终点则从头开始
	player.Position = (player.Position + steps) % 40

	// 经过起点奖励
	if player.Position < steps {
		player.Money += 200
	}

	return nil
}

// 当前玩家购买地产
func (g *Game) BuyProperty() error {
	player := g.Players[g.Current]
	space := g.Board[player.Position]

	if space == nil {
		return errors.New("当前位置不可购买")
	}

	if space.Owner != nil {
		return errors.New("该地产已被购买")
	}

	if player.Money < space.Price {
		return errors.New("资金不足")
	}

	// 购买地产
	player.Money -= space.Price
	space.Owner = player
	player.Property = append(player.Property, space)

	return nil
}

// 支付租金
func (g *Game) PayRent() error {
	player := g.Players[g.Current]
	space := g.Board[player.Position]

	if space == nil || space.Owner == nil || space.Owner == player {
		return nil // 不需要支付租金
	}

	if player.Money < space.Rent {
		return errors.New("资金不足以支付租金")
	}

	// 支付租金
	player.Money -= space.Rent
	space.Owner.Money += space.Rent

	return nil
}

// 切换到下一个玩家
func (g *Game) NextPlayer() {
	g.Current = (g.Current + 1) % len(g.Players)
}

// 检查游戏是否结束
func (g *Game) IsGameOver() bool {
	for _, player := range g.Players {
		if player.Money < 0 {
			return true
		}
	}
	return false
}

// 获取获胜者
func (g *Game) GetWinner() *Player {
	winner := g.Players[0]
	for _, player := range g.Players[1:] {
		if player.Money > winner.Money {
			winner = player
		}
	}
	return winner
}

// GenerateGameID 生成游戏ID
// TODO: 将改成链上VRF随机数实现
func GenerateGameID() string {
	// 实现一个简单的ID生成逻辑，实际项目中可以使用UUID
	return "game-" + string(rand.Int63())
}
