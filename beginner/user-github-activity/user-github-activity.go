package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

var cache = &sync.Map{}

// JSON unpack
type Actor struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

type Repo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PayLoad struct {
	Action string `json:"action"`
}

type Event struct {
	Id        string  `json:"id"`
	Type      string  `json:"type"`
	Actor     Actor   `json:"actor"`
	Repo      Repo    `json:"repo"`
	Payload   PayLoad `json:"payload"`
	Public    bool    `json:"public"`
	CreatedAt string  `json:"created_at"`
}

const URL = "https://api.github.com/users/"

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

func Menu() {
	CallClear()
	fmt.Printf("Enter username (enter exit to quit): ")
}

func GetEvent(username string) ([]Event, int) {
	var events []Event

	if result, ok := cache.Load(username); ok {
		return result.([]Event), 0
	}

	rawResult, err := http.Get(URL + username + "/events")
	if err != nil {
		return events, 1
	}
	defer rawResult.Body.Close()
	errJson := json.NewDecoder(rawResult.Body).Decode(&events)
	if errJson != nil {
		return events, 2
	}
	cache.Store(username, events)
	go func() {
		time.Sleep(10 * time.Minute)

		cache.Delete(username)
	}()

	return events, 0
}

func main() {
	ok := true
	scanner := bufio.NewScanner(os.Stdin)
	for ok {
		Menu()
		scanner.Scan()
		usrChoice := scanner.Text()

		if usrChoice == "exit" {
			ok = false
			continue
		}

		events, err := GetEvent(usrChoice)
		switch err {
		case 0:
			if len(events) == 0 {
				fmt.Println("Seems like the username doesn't exists, try again later. Press enter to continue...")
				scanner.Scan()
				continue
			}
			for i, event := range events {
				time, err := time.Parse("2006-01-02T15:04:05Z", event.CreatedAt)
				if err != nil {
					fmt.Println(err)
					return
				}
				if len(event.Payload.Action) == 0 {
					fmt.Printf("%d. %s - repo %s at: %s\n", i+1, event.Type[:len(event.Type)-5], event.Repo.Name, time.Format("Jan 02, 15:04"))
				} else {
					fmt.Printf("%d. %s - %s - repo %s at: %s\n", i+1, event.Type[:len(event.Type)-5], event.Payload.Action, event.Repo.Name, time.Format("Jan 02, 15:04"))
				}
			}
			fmt.Println("Successful! Press enter to continue...")
			scanner.Scan()
		case 1:
			fmt.Println("Couldn't get the data, please enter valid username. Press enter to continue...")
			scanner.Scan()
		case 2:
			fmt.Println("Something went wrong, please try agian. Press enter to continue...")
			scanner.Scan()
		}
	}
}
