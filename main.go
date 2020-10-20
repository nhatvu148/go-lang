package main

import "fmt"

// cmd: go run main.go deck.go
func main() {
	// var card string = "Ace of Spades"
	card := newCard()
	fmt.Println(card)

	// Array: static array, Slice: dynamic array
	// cards := deck{"Ace of Diamonds", newCard()}
	// cards = append(cards, "Six of Spades")
	// fmt.Println(cards)

	// cards.print()

	cards := newDeck()

	cards.print()
}

func newCard() string {
	return "Five of Hearts"
}
