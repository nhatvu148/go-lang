package main

func main() {
	// greeting := "Hello world!"
	// fmt.Println([]byte(greeting))

	// cards := newDeck()
	// fmt.Println(cards.toString())
	// cards.saveToFile("my-cards")

	newCards := newDeckFromFile("my-cards")
	newCards.print()
}
