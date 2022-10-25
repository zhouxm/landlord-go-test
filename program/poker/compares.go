package poker

// IsFirstCardValueBigger
//常规比较，按照值大小比较,如果第一个牌比第二个牌的值大，返回true
func IsFirstCardValueBigger(card1 *Card, card2 *Card) bool {
	if card1.GetValue() > card2.GetValue() {
		return true
	} else {
		return false
	}
}

func GetBigCard(card1 *Card, card2 *Card) *Card {
	if IsFirstCardValueBigger(card1, card2) {
		return card1
	} else {
		return card2
	}
}

func GetSmallCard(card1 *Card, card2 *Card) *Card {
	if IsFirstCardValueBigger(card1, card2) {
		return card2
	} else {
		return card1
	}
}
