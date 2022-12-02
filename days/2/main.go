package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Move string
type Outcome string

const (
  Rock Move = "rock"
  Paper = "paper"
  Scissors = "scissors"
)

const (
  Loss Outcome = "loss"
  Tie = "tie"
  Win = "win"
)

type Round struct {
  opponentMove Move
  myMove Move
  outcome Outcome
}

type Strategy struct {
  rounds []Round
  score int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := readInputFile(os.Args[1])
  strategy := parseStrategy(input)
fmt.Println("strategy score", strategy.score);
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseStrategy(input string) *Strategy {
  
  lines := strings.Split(input, "\n")

  rounds := make([]Round, len(lines))
  score := 0

  for i, roundInput := range lines {
    if len(roundInput) == 0 { continue }
    rounds[i] = parseRound2(roundInput)   
    score += computeRoundScore(rounds[i])
  }
  
 
  return &Strategy{
    rounds: rounds,
    score: score,
  }
}

func parseRound(input string) Round {
  oMove := parseOpponentMove(input[0])
  mMove := parseMyMove(input[2])
  round := Round{ 
    opponentMove: oMove, 
    myMove: mMove,
    outcome: getOutcome(oMove, mMove),
  }
  fmt.Println(round)
  return round
}

func parseRound2(input string) Round {
  oMove := parseOpponentMove(input[0])
  outcome := parseMyOutcome(input[2])
  mMove := determineMyMove(outcome, oMove)
  round := Round{ 
    opponentMove: oMove, 
    myMove: mMove,
    outcome: getOutcome(oMove, mMove),
  }
  fmt.Println(round)
  return round
}

func computeRoundScore(round Round) int {
  var outcomeScore, choiceScore int
  switch round.outcome {
    case Win: outcomeScore = 6
    case Tie: outcomeScore = 3
    case Loss: outcomeScore = 0
    default: panic(errors.New("Unknown outcome type"))
  }

  switch round.myMove {
    case Rock: choiceScore = 1
    case Paper: choiceScore = 2
    case Scissors: choiceScore = 3
    default: panic(errors.New("Unknown outcome type"))
  }

  return outcomeScore + choiceScore 
}

func parseOpponentMove(input byte) Move {
  if input == 'A' { return Rock }
  if input == 'B' { return Paper }
  if input == 'C' { return Scissors }
  panic(errors.New("Unknown move type"))
}

func parseMyMove(input byte) Move {
  if input == 'X' { return Rock }
  if input == 'Y' { return Paper }
  if input == 'Z' { return Scissors }
  panic(errors.New("Unknown move type"))
}

func getOutcome(opponentMove Move, myMove Move) Outcome {
  if opponentMove == Rock && myMove == Paper { return Win }
  if opponentMove == Rock && myMove == Scissors { return Loss }
 
  if opponentMove == Paper && myMove == Scissors { return Win }
  if opponentMove == Paper && myMove == Rock { return Loss }
 
  if opponentMove == Scissors && myMove == Rock { return Win }
  if opponentMove == Scissors && myMove == Paper { return Loss }

  return Tie
}

func parseMyOutcome(input byte) Outcome {
  if input == 'X' { return Loss }
  if input == 'Y' { return Tie }
  if input == 'Z' { return Win }
  panic(errors.New("Unknown move type"))
}

func determineMyMove(outcome Outcome, opponentMove Move) Move {
  if outcome == Loss {
    switch opponentMove {
      case Rock: return Scissors 
      case Paper: return Rock 
      case Scissors: return Paper 
      default: panic(errors.New("Unknown move type"))
    }
  } else if outcome == Win {
    switch opponentMove {
      case Rock: return Paper 
      case Paper: return Scissors 
      case Scissors: return Rock 
      default: panic(errors.New("Unknown move type"))
    }
  } else {
    return opponentMove
  }
}
