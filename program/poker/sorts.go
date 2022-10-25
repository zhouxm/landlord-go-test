package poker

// BubbleSortCardsMax2Min 使用冒泡排序法，对给定的扑克牌，使用给定的规则进项从小到大排序
func BubbleSortCardsMax2Min(cards CardSet, maxCard func(card1 *Card, card2 *Card) bool) {
	length := cards.CountCards()
	for i := 0; i < length; i++ {
		for j := i; j < length; j++ {
			if !maxCard(cards[i], cards[j]) {
				cards[i], cards[j] = cards[j], cards[i]
			}
		}
	}
}

// BubbleSortCardsMin2Max 使用冒泡排序法，对给定的扑克牌，使用给定的规则进项从小到大排序
func BubbleSortCardsMin2Max(cards CardSet, maxCard func(card1 *Card, card2 *Card) bool) {
	length := cards.CountCards()
	for i := 0; i < length; i++ {
		for j := i; j < length; j++ {
			if maxCard(cards[i], cards[j]) {
				cards[i], cards[j] = cards[j], cards[i]
			}
		}
	}
}
