package games

import (
	"landlord/program/game"
)

/**
*该包用于解决game和doudizhu包循环依赖问题
 */
func NewGame(gameID int, baseScore int) game.IGame {
	switch gameID {
	case game.TypeOfDoudozhu:
		return GetLandlord(baseScore)
	case game.TypeOfShengji:
		return nil
	case game.TypeOfBaohuang:
		return nil
	case game.TypeOfZhajinhua:
		return nil
	default:
		return nil
	}
}
