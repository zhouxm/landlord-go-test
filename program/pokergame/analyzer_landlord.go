package pokergame

import (
	"landlord/program/poker"
	"sync"
)

//定义玩家的扑克牌分析器map的索引为poker的value,value为改值得扑克牌在玩家牌中的索引
type landLordAnalyzer struct {
	sync.RWMutex
	dic map[int]poker.CardSet
}

//根据给定的扑克集初始化分析器
func (ana *landLordAnalyzer) InitAnalyzer() {
	ana.dic = make(map[int]poker.CardSet)
	ana.dic[poker.CARD_VALUE_THREE] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_FOUR] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_FIVE] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_SIX] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_SEVEN] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_EIGHT] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_NINE] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_TEN] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_JACK] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_QUEEN] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_KING] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_ACE] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_TWO] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_BLACK_JOKER] = poker.CardSet{}
	ana.dic[poker.CARD_VALUE_RED_JOKER] = poker.CardSet{}
}

//根据给定的扑克集更新记牌器,出牌时调用
func (ana *landLordAnalyzer) RemovePokerSet(pokers poker.CardSet) {
	ana.Lock()
	defer ana.Unlock()
	pokers.DoOnEachPokerCard(func(index int, card *poker.Card) {
		ana.dic[card.GetValue()], _ = ana.dic[card.GetValue()].DelPokers(poker.CardSet{card})
	})
}

func (ana *landLordAnalyzer) AddPokerSet(pokers poker.CardSet) {
	ana.Lock()
	defer ana.Unlock()
	pokers.DoOnEachPokerCard(func(index int, card *poker.Card) {
		ana.dic[card.GetValue()] = ana.dic[card.GetValue()].AddPokers(poker.CardSet{card})
	})
}

func (ana *landLordAnalyzer) GetMinPlayableCards() poker.CardSet {
	ana.Lock()
	defer ana.Unlock()
	for i := poker.CARD_VALUE_THREE; i <= poker.CARD_VALUE_RED_JOKER; i++ {
		set, _ := ana.dic[i]
		if set.CountCards() > 0 {
			return set
		}
	}
	return poker.CardSet{}
}

//根据最后一次出牌的牌型信息，返回可出的扑克集
func (ana *landLordAnalyzer) GetUseableCards(setType *SetInfo) []poker.CardSet {
	ana.Lock()
	defer ana.Unlock()

	var useableSets []poker.CardSet

	switch setType.setType {
	case LANDLORD_SET_TYPE_SINGLE:
		useableSets = ana.getSingleValueSet(1, setType.GetMinValue())
	case LANDLORD_SET_TYPE_DRAGON:
		useableSets = ana.getMultiValueSet(1, setType.GetMinValue(), setType.GetMaxValue())
	case LANDLORD_SET_TYPE_PAIR:
		useableSets = ana.getSingleValueSet(2, setType.GetMinValue())
	case LANDLORD_SET_TYPE_MULIT_PAIRS:
		useableSets = ana.getMultiValueSet(2, setType.GetMinValue(), setType.GetMaxValue())
	case LANDLORD_SET_TYPE_THREE:
		useableSets = ana.getSingleValueSet(3, setType.GetMinValue())
	case LANDLORD_SET_TYPE_THREE_PLUS_ONE:
		useableSets = ana.getSingleValueSet(3, setType.GetMinValue())
		for i, tempset := range useableSets {
			tempsetPlus := ana.getPlusSet(1, 1, tempset)
			if tempsetPlus.CountCards() > 0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			} else { //没有牌可以带，将之前的主牌移除可出牌集合
				useableSets[i] = nil
			}
		}
	case LANDLORD_SET_TYPE_THREE_PLUS_TWO:
		useableSets = ana.getSingleValueSet(3, setType.GetMinValue())
		for i, tempset := range useableSets {
			tempsetPlus := ana.getPlusSet(2, 1, tempset)
			if tempsetPlus.CountCards() > 0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			} else { //没有牌可以带，将之前的主牌移除可出牌集合
				useableSets[i] = nil
			}
		}
	case LANDLORD_SET_TYPE_MULITY_THREE:
		useableSets = ana.getMultiValueSet(3, setType.GetMinValue(), setType.GetMaxValue())
	case LANDLORD_SET_TYPE_MULITY_THREE_PLUS_ONE:
		useableSets = ana.getMultiValueSet(3, setType.GetMinValue(), setType.GetMaxValue())
		for i, tempset := range useableSets {
			tempsetPlus := ana.getPlusSet(1, setType.GetRangeWidth(), tempset)
			if tempsetPlus.CountCards() > 0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			} else { //没有牌可以带，将之前的主牌移除可出牌集合
				useableSets[i] = nil
			}
		}
	case LANDLORD_SET_TYPE_MULITY_THREE_PLUS_TWO:
		useableSets = ana.getMultiValueSet(3, setType.GetMinValue(), setType.GetMaxValue())
		for i, tempset := range useableSets {
			tempsetPlus := ana.getPlusSet(2, setType.GetRangeWidth(), tempset)
			if tempsetPlus.CountCards() > 0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			} else { //没有牌可以带，将之前的主牌移除可出牌集合
				useableSets[i] = nil
			}
		}
	case LANDLORD_SET_TYPE_FOUR_PLUS_TWO:
		useableSets = ana.getSingleValueSet(4, setType.GetMinValue())
		for i, tempset := range useableSets {
			//带两个单牌
			tempsetPlus := ana.getPlusSet(1, 2, tempset)
			if tempsetPlus.CountCards() > 0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			} else {
				//带一对牌，看做两个单牌
				tempsetPlus := ana.getPlusSet(2, 1, tempset)
				if tempsetPlus.CountCards() > 0 {
					useableSets[i] = tempset.AddPokers(tempsetPlus)
				} else { //没有牌可以带，将之前的主牌移除可出牌集合
					useableSets[i] = nil
				}
			}
		}
	case LANDLORD_SET_TYPE_FOUR_PLUS_FOUR:
		useableSets = ana.getSingleValueSet(4, setType.GetMinValue())
		for i, tempset := range useableSets {
			tempsetPlus := ana.getPlusSet(2, 2, tempset)
			if tempsetPlus.CountCards() > 0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			} else { //没有牌可以带，将之前的主牌移除可出牌集合
				useableSets[i] = nil
			}
		}
	case LANDLORD_SET_TYPE_MULITY_FOUR:
		useableSets = ana.getMultiValueSet(4, setType.GetMinValue(), setType.GetMaxValue())
	case LANDLORD_SET_TYPE_MULITY_FOUR_PLUS_TWO:
		useableSets = ana.getMultiValueSet(4, setType.GetMinValue(), setType.GetMaxValue())
		for i, tempset := range useableSets {
			//带两个单牌
			tempsetPlus := ana.getPlusSet(1, 2*setType.GetRangeWidth(), tempset)
			if tempsetPlus.CountCards() > 0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			} else {
				//带一对牌，看做两个单牌
				tempsetPlus := ana.getPlusSet(2, setType.GetRangeWidth(), tempset)
				if tempsetPlus.CountCards() > 0 {
					useableSets[i] = tempset.AddPokers(tempsetPlus)
				} else { //没有牌可以带，将之前的主牌移除可出牌集合
					useableSets[i] = nil
				}
			}
		}
	case LANDLORD_SET_TYPE_MULITY_FOUR_PLUS_FOUR:
		useableSets = ana.getMultiValueSet(4, setType.GetMinValue(), setType.GetMaxValue())
		for i, tempset := range useableSets {
			tempsetPlus := ana.getPlusSet(2, 2*setType.GetRangeWidth(), tempset)
			if tempsetPlus.CountCards() > 0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			} else { //没有牌可以带，将之前的主牌移除可出牌集合
				useableSets[i] = nil
			}
		}
	case LANDLORD_SET_TYPE_COMMON_BOMB:
		useableSets = ana.getSingleValueSet(4, setType.GetMinValue())
	case LANDLORD_SET_TYPE_JOKER_BOMB:
		useableSets = []poker.CardSet{}
	default:
		useableSets = []poker.CardSet{}
	}
	//去掉nil元素
	newUseableSets := []poker.CardSet{}
	for _, sets := range useableSets {
		if sets != nil {
			newUseableSets = append(newUseableSets, sets)
		}
	}
	//上一次出牌不是炸弹，则直接将炸弹加入可出的排中
	if setType.setType != LANDLORD_SET_TYPE_COMMON_BOMB && setType.setType != LANDLORD_SET_TYPE_JOKER_BOMB {
		//王炸
		jokerBombSet := ana.getJokerBomb()
		if jokerBombSet.CountCards() > 0 {
			newUseableSets = append(newUseableSets, jokerBombSet)
		}
		//普通炸弹
		for _, tempSet := range ana.getSingleValueSet(4, -1) {
			if tempSet.CountCards() > 0 {
				newUseableSets = append(newUseableSets, tempSet)
			}
		}
	}
	return newUseableSets
}

//获取单值牌组成的扑克集的切片，单排对牌三牌四排等等
//count表示单值牌的张数
//minValue表示上家出牌的最小的牌的大小
func (ana *landLordAnalyzer) getSingleValueSet(count int, minValue int) []poker.CardSet {
	sets := []poker.CardSet{}
	se := poker.NewPokerSet()
	//先不拆牌的情况下查找
	for i := minValue + 1; i <= poker.CARD_VALUE_RED_JOKER; i++ {
		if ana.dic[i].CountCards() == count {
			se = se.AddPokers(ana.dic[i])
			sets = append(sets, se)
			se = poker.NewPokerSet()
		}
	}
	//不拆牌的情况下找不到可出的牌，再考虑拆牌的情况
	if len(sets) == 0 {
		for i := minValue + 1; i <= poker.CARD_VALUE_RED_JOKER; i++ {
			if ana.dic[i].CountCards() > count {
				se = se.AddPokers(ana.dic[i][:count])
				sets = append(sets, se)
				se = poker.NewPokerSet()
			}
		}
	}
	return sets
}

//获取多种不同值组成的扑克集的切片,2连3连4连5连等
func (ana *landLordAnalyzer) getMultiValueSet(count int, minValue int, maxValue int) []poker.CardSet {
	sets := []poker.CardSet{}
	se := poker.NewPokerSet()
	valueRange := maxValue - minValue + 1
	//先考虑不拆拍的情况
	for i := minValue + 1; i <= poker.CARD_VALUE_TWO-valueRange; i++ {
		for j := i; j < i+valueRange; j++ {
			if ana.dic[j].CountCards() == count {
				se = se.AddPokers(ana.dic[j])
			}
		}
		//该范围内连续的牌的张数符合要求
		if se.CountCards() == valueRange*count {
			sets = append(sets, se)
			se = poker.NewPokerSet()
		} else {
			se = poker.NewPokerSet()
		}
	}
	//如果不拆拍找不到可出的牌，则考虑拆牌
	if len(sets) == 0 {
		for i := minValue + 1; i <= poker.CARD_VALUE_TWO-valueRange; i++ {
			for j := i; j < i+valueRange; j++ {
				if ana.dic[j].CountCards() > count {
					se = se.AddPokers(ana.dic[j][:count])
				}
			}
			//该范围内连续的牌的张数符合要求
			if se.CountCards() == valueRange*count {
				sets = append(sets, se)
				se = poker.NewPokerSet()
			} else {
				se = poker.NewPokerSet()
			}
		}
	}

	return sets
}

//获取附牌，比如三带一中的一，四带二中二，只获取一种可能即可
//不拆牌为第一原则，可能会带出去大牌
//num张数count系列数exceptset不能包含在内的扑克集
func (ana *landLordAnalyzer) getPlusSet(num int, count int, exceptSet poker.CardSet) poker.CardSet {
	resSet := poker.NewPokerSet()
	//第一原则不拆牌原则
	for i := poker.CARD_VALUE_THREE; i <= poker.CARD_VALUE_RED_JOKER; i++ {
		if ana.dic[i].CountCards() == num {
			if !ana.dic[i][:num].HasSameValueCard(exceptSet) {
				resSet = resSet.AddPokers(ana.dic[i])
			}
		}
		if resSet.CountCards() == num*count {
			return resSet
		}
	}
	//不拆牌找不到则，考虑拆牌
	if resSet.CountCards() == 0 {
		for i := poker.CARD_VALUE_THREE; i <= poker.CARD_VALUE_RED_JOKER; i++ {
			if ana.dic[i].CountCards() > num {
				if !ana.dic[i][:num].HasSameValueCard(exceptSet) {
					resSet = resSet.AddPokers(ana.dic[i][:num])
				}
			}
			if resSet.CountCards() == num*count {
				return resSet
			}
		}
	}

	return poker.CardSet{}
}
func (ana *landLordAnalyzer) getJokerBomb() poker.CardSet {
	resSet := poker.NewPokerSet()
	for i := poker.CARD_VALUE_BLACK_JOKER; i <= poker.CARD_VALUE_RED_JOKER; i++ {
		if ana.dic[i].CountCards() > 0 {
			resSet = resSet.AddPokers(ana.dic[i])
		}
	}
	if resSet.CountCards() > 1 {
		return resSet
	} else {
		return poker.NewPokerSet()
	}
}
