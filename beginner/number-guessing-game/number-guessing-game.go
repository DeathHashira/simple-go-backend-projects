package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

var bestTry int

// clear CLI
var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func HowFar(guess int, goal int) int {
	distance := int(math.Abs(float64(guess - goal)))
	if goal > guess {
		if distance == 1 {
			return 0
		} else if distance < 5 {
			return 1
		} else if distance < 17 {
			return 2
		} else if distance < 65 {
			return 3
		} else {
			return 4
		}

	} else {
		if distance == 1 {
			return 5
		} else if distance < 5 {
			return 6
		} else if distance < 17 {
			return 7
		} else if distance < 65 {
			return 8
		} else {
			return 9
		}
	}
}

func Menu() {
	fmt.Printf("Welcome to the Number Guessing Game!\nI'm thinking of a number between 1 and 100.\n\nPlease select the difficulty level:\n1. Easy (10 chances)\n2. Medium (5 chances)\n3. Hard (3 chances)\n")
}

func Count(attemp int) {
	fmt.Printf("You have %d chance left!\n", attemp)
}

func main() {

	CallClear()
	scanner := bufio.NewScanner(os.Stdin)
	redo := true
	var chanceNum int

	for redo {
		won := false
		attemps := 1
		Menu()
		fmt.Printf("Enter your choice: ")
		scanner.Scan()
		difficulty, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("It's not a valid choice. Do it again. Press enter to continue...")
			scanner.Scan()
			CallClear()
			continue
		}
		switch difficulty {
		case 1:
			chanceNum = 10
		case 2:
			chanceNum = 5
		case 3:
			chanceNum = 3
		default:
			println("It's not a valid choice. Do it again. Press enter to continue...")
			scanner.Scan()
			CallClear()
			continue
		}
		guessingNum := rand.Intn(99) + 1
		CallClear()
		Count(chanceNum - attemps + 1)

		for i := 0; i < chanceNum; i++ {
			fmt.Printf("Enter your guess: ")
			scanner.Scan()
			usrNum, err := strconv.Atoi(scanner.Text())
			CallClear()
			Count(chanceNum - attemps)
			if err != nil {
				fmt.Println("Oh no. You just lost your attempt with invalid choice. Be more carefull.")
				attemps += 1
				continue
			}
			if guessingNum == usrNum {
				won = true
				break
			} else {
				if chanceNum-attemps == 0 {
					CallClear()
					break
				}
				switch HowFar(usrNum, guessingNum) {
				case 0:
					fmt.Println("Wow. You're almost guessed it right! But you have to guess a little bit up.")
				case 1:
					fmt.Println("You're so close. But your guess is still less than the number.")
				case 2:
					fmt.Println("You're now up to the number. You have to guess higher.")
				case 3:
					fmt.Println("Hmm! Not close. You have to guess a higher numebr.")
				case 4:
					fmt.Println("Oops! You're too far. Guess more higher number.")
				case 5:
					fmt.Println("You're sooo close! Just a little bit less.")
				case 6:
					fmt.Println("You're getting close. But your guess is still upper than the number.")
				case 7:
					fmt.Println("You're in right direction. But have to guess a low number.")
				case 8:
					fmt.Println("Hmm! Not so close. You have to guess a lower number.")
				case 9:
					fmt.Println("Oops! You're too far. Guess a much lower number.")
				}
				attemps += 1
			}
		}
		if won {
			CallClear()
			if bestTry == 0 {
				bestTry = attemps
			} else {
				bestTry = int(math.Min(float64(bestTry), float64(attemps)))
			}
			fmt.Printf("Congratulations! You guessed the correct number in %d attempts.\nbest record: %d attempt\n", attemps, bestTry)
		} else {
			fmt.Printf("Oops! You lost :(.\n")
		}
		fmt.Println("Do you wanna play again?(Y/N): ")
		scanner.Scan()
		usrChoice := scanner.Text()
		switch usrChoice {
		case "y", "Y":
			redo = true
		case "N", "n":
			redo = false
		default:
			fmt.Printf("Invalid option. Please try again later. Press enter to continue...")
			scanner.Scan()
			return
		}
		CallClear()
	}
}
