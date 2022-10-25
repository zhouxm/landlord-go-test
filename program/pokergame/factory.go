package pokergame

import (
	"landlord/program/poker"
)

//定义扑克牌游戏的类型，用于获取适用于不同扑克牌的set和analayzer
const (
	GAME_OF_LANDLORD = iota
	GAME_OF_ROYALIST
	GAME_OF_ZHAJINHUA
	GAME_OF_SHENGJI
)

type IRecorder interface {
	InitRecorder()
	AddPokerSet(playerPokers ...poker.CardSet)
	RemovePokerSet(cards poker.CardSet)
	SequenceJsonEncode() string
}

type IAnalyzer interface {
	InitAnalyzer()
	AddPokerSet(pokers poker.CardSet)
	RemovePokerSet(pokers poker.CardSet)
	GetMinPlayableCards() poker.CardSet
	GetUseableCards(setType *SetInfo) []poker.CardSet
}

type ISetChecker interface {
	GetSetInfo(set poker.CardSet) (*SetInfo, error)
}

//获取适用于某一种游戏的扑克集
func NewSetChecker(gameType int) ISetChecker {
	switch gameType {
	case GAME_OF_LANDLORD:
		landLordChecker := landLordChecker{}
		return &landLordChecker
	case GAME_OF_ROYALIST:
	case GAME_OF_SHENGJI:
	case GAME_OF_ZHAJINHUA:
	default:
		return nil
	}
	return nil
}

//获取适用于某一种游戏的扑克分析器
func NewAnalyzer(gameType int) IAnalyzer {
	switch gameType {
	case GAME_OF_LANDLORD:
		analyzer := landLordAnalyzer{}
		analyzer.InitAnalyzer()
		return &analyzer
	case GAME_OF_ROYALIST:
	case GAME_OF_SHENGJI:
	case GAME_OF_ZHAJINHUA:
	default:
		return nil
	}
	return nil
}

//获取适用于某一种游戏的扑克记录器
func NewRecorder(gameType int) IRecorder {
	switch gameType {
	case GAME_OF_LANDLORD:
		recorder := cardNameRecorder{}
		recorder.InitRecorder()
		return &recorder
	case GAME_OF_ROYALIST:
	case GAME_OF_SHENGJI:
	case GAME_OF_ZHAJINHUA:
	default:
		return nil
	}
	return nil
}
