package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Card int

const (
	CardJoker Card = iota
	Card2
	Card3
	Card4
	Card5
	Card6
	Card7
	Card8
	Card9
	Card10
	CardQueen
	CardKing
	CardAce
)

type Rank int

const (
	RankHighCard Rank = iota
	RankOnePair
	RankTwoPair
	RankThreeOfAKind
	RankFullHouse
	RankFourOfAKind
	RankFiveOfAKind
)

type Hand struct {
	Bid  int
	Hand [5]Card
	Rank Rank
}

func runMain() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	var hands []Hand
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		res, err := parseHand(line)
		if err != nil {
			return fmt.Errorf("parse hand: %w", err)
		}

		hands = append(hands, res)
	}

	slices.SortFunc(hands, func(a, b Hand) int {
		if n := cmp.Compare(a.Rank, b.Rank); n != 0 {
			return n
		}

		return slices.Compare(a.Hand[:], b.Hand[:])
	})

	winnings := 0
	for i, hand := range hands {
		winnings += hand.Bid * (i + 1)
	}

	fmt.Println(winnings)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	return nil
}

func parseHand(text string) (Hand, error) {
	parts := strings.SplitN(text, " ", 2)

	var hand [5]Card
	for i, card := range parts[0] {
		if card >= '2' && card <= '9' {
			hand[i] = Card(card - '1')
			continue
		}

		switch card {
		case 'J':
			hand[i] = CardJoker
		case 'T':
			hand[i] = Card10
		case 'Q':
			hand[i] = CardQueen
		case 'K':
			hand[i] = CardKing
		case 'A':
			hand[i] = CardAce
		default:
			return Hand{}, fmt.Errorf("unknown card %q", hand[i])
		}
	}

	bid, err := strconv.Atoi(parts[1])
	if err != nil {
		return Hand{}, fmt.Errorf("parse bid %q: %w", parts[1], err)
	}

	return Hand{
		Bid:  bid,
		Hand: hand,
		Rank: rankHand(hand),
	}, nil
}

func rankHand(hand [5]Card) Rank {
	cardCount := make(map[Card]int, 5)
	maxCount := 0
	jokerCount := 0

	for _, card := range hand {
		if card == CardJoker {
			jokerCount++
			continue
		}

		cardCount[card]++

		if cardCount[card] > maxCount {
			maxCount = cardCount[card]
		}
	}

	if maxCount == 0 {
		return RankFiveOfAKind // all jokers
	}

	if maxCount == 1 {
		switch jokerCount {
		case 0:
			return RankHighCard
		case 1:
			return RankOnePair
		case 2:
			return RankThreeOfAKind
		case 3:
			return RankFourOfAKind
		case 4:
			return RankFiveOfAKind
		}
	}

	if maxCount == 2 {
		pairs := 0
		for _, count := range cardCount {
			if count == 2 {
				pairs++
			}
		}

		if pairs == 2 {
			switch jokerCount {
			case 0:
				return RankTwoPair
			case 1:
				return RankFullHouse
			}
		} else {
			switch jokerCount {
			case 0:
				return RankOnePair
			case 1:
				return RankThreeOfAKind
			case 2:
				return RankFourOfAKind
			case 3:
				return RankFiveOfAKind
			}
		}
	}

	if maxCount == 3 {
		hasPair := false
		for _, count := range cardCount {
			if count == 2 {
				hasPair = true
				break
			}
		}

		if hasPair {
			return RankFullHouse
		}

		switch jokerCount {
		case 0:
			return RankThreeOfAKind
		case 1:
			return RankFourOfAKind
		case 2:
			return RankFiveOfAKind
		}
	}

	if maxCount == 4 {
		return Rank(int(RankFourOfAKind) + jokerCount)
	}

	if maxCount == 5 {
		return RankFiveOfAKind
	}

	panic(fmt.Errorf("unexpected rank for hand %q", hand))
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
