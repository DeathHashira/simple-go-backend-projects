package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// create user-side cache
var cache = &sync.Map{}

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

// structure type for JSON response
type Hour struct {
	Datetime  string  `json:"datetime"`
	Temp      float64 `json:"temp"`
	Feelslike float64 `json:"feelslike"`
	Humidity  float64 `json:"humidity"`
	Icon      string  `json:"icon"`
}

type Day struct {
	Datetime  string  `json:"datetime"`
	Humidity  float64 `json:"humidity"`
	Temp      float64 `json:"temp"`
	Icon      string  `json:"icon"`
	Feelslike float64 `json:"feelslike"`
	Hours     []Hour  `json:"hours"`
}

type Weather struct {
	Timezone string `json:"timezone"`
	Address  string `json:"address"`
	Days     []Day  `json:"days"`
}

// API URL
const baseURL = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/"

// unique name for each data
func CacheHash(city string, dates []string) string {
	if len(dates) == 0 {
		now := time.Now().Format("2006-01-02")
		next := time.Now().AddDate(0, 0, 15).Format("2006-01-02")
		return city + "|" + now + "|" + next
	} else if len(dates) == 1 {
		return city + "|" + dates[0]
	} else {
		return city + "|" + dates[0] + "|" + dates[1]
	}
}

// API get request and unpack JSON
func GetWeather(city string, dates []string) (Weather, int) {
	var weather Weather
	var url string

	entity := CacheHash(city, dates)
	if result, ok := cache.Load(entity); ok {
		return result.(Weather), 0
	}

	if len(dates) == 1 {
		url = baseURL + city + "/" + dates[0] + "?key=" + os.Getenv("API_KEY")
	} else if len(dates) == 0 {
		url = baseURL + city + "?key=" + os.Getenv("API_KEY")
	} else {
		url = baseURL + city + "/" + dates[0] + "/" + dates[1] + "?key=" + os.Getenv("API_KEY")
	}

	resp, err := http.Get(url)
	if err != nil {
		return weather, -1
	}
	if resp.StatusCode != http.StatusOK {
		return weather, -1
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil && err != io.EOF {
		return weather, -1
	}

	cache.Store(entity, weather)
	go func() {
		time.Sleep(10 * time.Minute)
		cache.Delete(entity)
	}()

	return weather, 0
}

func Menu() {
	CallClear()
	fmt.Printf("What you want to do: \n1. Weather on current time\n2. Weather for today\n3. Weather for next 15 days\n4. Weather on specific time\n5. Exit\n")
}

func main() {
	// get the API key
	godotenv.Load()
	scanner := bufio.NewScanner(os.Stdin)

	var usrChoice int
	for ok := true; ok; ok = (usrChoice != 5) {
		Menu()
		fmt.Printf("Enter your choice: ")
		scanner.Scan()
		usrChoice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			fmt.Printf("Invalid choice. Please try again. Press enter to continue...")
			scanner.Scan()
			continue
		}
		if usrChoice == 5 {
			break
		} else if usrChoice > 5 {
			fmt.Printf("Invalid choice. Please try again. Press enter to continue...")
			scanner.Scan()
			continue
		}
		fmt.Printf("What country/city you want to get the weather of: ")
		scanner.Scan()
		timeZone := scanner.Text()
		switch usrChoice {
		case 1:
			date := time.Now().Format("2006-01-02")
			hourShow := time.Now().Format("15") + ":00:00"
			dates := [1]string{date}
			weather, err := GetWeather(timeZone, dates[:])
			if err == -1 {
				fmt.Println("Couldn't get the data, please try again. Press enter to continue...")
				scanner.Scan()
				continue
			}
			t1, err1 := time.Parse("2006-01-02", weather.Days[0].Datetime)
			if err1 != nil {
				fmt.Println(err1)
				return
			}
			CallClear()
			fmt.Printf("%s - %s - %s\n\n", weather.Timezone, weather.Address, t1.Format("Mon, Jan 2"))
			for _, d := range weather.Days {
				for _, h := range d.Hours {
					if h.Datetime == hourShow {
						fmt.Printf("      %s      Temperature: %0.1f\tHumidity: %0.1f   \tFeels like: %0.1f \t %s", hourShow[:len(hourShow)-3], h.Temp, h.Humidity, h.Feelslike, h.Icon)
					}
				}
			}
			fmt.Print("\nPress enter to go back to menu...")
			scanner.Scan()

		case 2:
			date := time.Now().Format("2006-01-02")
			dates := [1]string{date}
			weather, err := GetWeather(timeZone, dates[:])
			if err == -1 {
				fmt.Println("Couldn't get the data, please try again. Press enter to continue...")
				scanner.Scan()
				continue
			}
			t1, err1 := time.Parse("2006-01-02", weather.Days[0].Datetime)
			if err1 != nil {
				fmt.Println(err1)
				return
			}
			CallClear()
			fmt.Printf("%s - %s - %s\n\n", weather.Timezone, weather.Address, t1.Format("Mon, Jan 2"))
			for _, d := range weather.Days {
				for _, h := range d.Hours {
					for j := 0; j < 5; j++ {
						switch j {
						case 0:
							fmt.Printf("      %s      ", h.Datetime[:5])
						case 1:
							fmt.Printf("Temperature: %0.1f\t", h.Temp)
						case 2:
							fmt.Printf("Humidity: %0.1f   \t", h.Humidity)
						case 3:
							fmt.Printf("Feels like: %0.1f \t", h.Feelslike)
						case 4:
							fmt.Printf("%s\t", h.Icon)
						}
					}
					fmt.Printf("\n")
				}
			}
			fmt.Print("\nPress enter to go back to menu...")
			scanner.Scan()

		case 3:
			dates := [0]string{}
			weather, err := GetWeather(timeZone, dates[:])
			if err == -1 {
				fmt.Println("Couldn't get the data, please try again. Press enter to continue...")
				scanner.Scan()
				continue
			}
			CallClear()
			fmt.Printf("%s - %s\n\n", weather.Timezone, weather.Address)
			for _, d := range weather.Days {
				t1, err := time.Parse("2006-01-02", d.Datetime)
				if err != nil {
					fmt.Println(err)
					return
				}
				for j := 0; j < 5; j++ {
					switch j {
					case 0:
						fmt.Printf("   %s   ", t1.Format("Mon, Jan 02"))
					case 1:
						fmt.Printf("Temperature: %0.1f\t", d.Temp)
					case 2:
						fmt.Printf("Humidity: %0.1f   \t", d.Humidity)
					case 3:
						fmt.Printf("Feels like: %0.1f \t", d.Feelslike)
					case 4:
						fmt.Printf("%s\t", d.Icon)
					}
				}
				fmt.Printf("\n")
			}
			fmt.Print("\nPress enter to go back to menu...")
			scanner.Scan()

		case 4:
			fmt.Printf("What is your start date(format: YYYY-MM-DD): ")
			scanner.Scan()
			date1 := scanner.Text()
			fmt.Printf("What is your end date(format: YYYY-MM-DD): ")
			scanner.Scan()
			date2 := scanner.Text()
			dates := [2]string{date1, date2}
			weather, err := GetWeather(timeZone, dates[:])
			if err == -1 {
				fmt.Println("Couldn't get the data, please try again. Press enter to continue...")
				scanner.Scan()
				continue
			}
			CallClear()
			fmt.Printf("%s - %s\n\n", weather.Timezone, weather.Address)
			for _, d := range weather.Days {
				t1, err := time.Parse("2006-01-02", d.Datetime)
				if err != nil {
					fmt.Println(err)
					return
				}
				for j := 0; j < 5; j++ {
					switch j {
					case 0:
						fmt.Printf("   %s   ", t1.Format("Mon, Jan 02"))
					case 1:
						fmt.Printf("Temperature: %0.1f\t", d.Temp)
					case 2:
						fmt.Printf("Humidity: %0.1f   \t", d.Humidity)
					case 3:
						fmt.Printf("Feels like: %0.1f \t", d.Feelslike)
					case 4:
						fmt.Printf("%s\t", d.Icon)
					}
				}
				fmt.Printf("\n")
			}
			fmt.Print("\nPress enter to go back to menu...")
			scanner.Scan()

		default:
			fmt.Printf("Invalid input. Please try again. Press enter to continue...")
			scanner.Scan()
		}

	}
}
