package luhn

func IsValidCard(cardNumber int) bool {
	return (cardNumber%10+luhn(cardNumber/10))%10 == 0
}
func luhn(cardNumber int) int {
	var sum int
	for i := 0; cardNumber > 0; i++ {
		currentDigit := cardNumber % 10
		if i%2 == 0 {
			currentDigit = currentDigit * 2
			if currentDigit > 9 {
				currentDigit = currentDigit - 9
			}
		}
		sum += currentDigit
		cardNumber = cardNumber / 10
	}
	return sum % 10
}
