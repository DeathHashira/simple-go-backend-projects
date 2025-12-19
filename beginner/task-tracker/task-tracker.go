package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

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

// print menu
func menu() {
	fmt.Println("-------------------------------")
	fmt.Print("What you want to do:\n1. add task\n2. modify task\n3. delete task\n4. update task\n5. list tasks\n6. exit\n")
	fmt.Println("-------------------------------")
	fmt.Print("Enter your choice: ")
}

// add new task
func add_task(title string, status string) {
	f, err := os.OpenFile(
		"tracker.json",
		os.O_CREATE|os.O_RDONLY,
		0644,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	var tasks []Task
	err = json.NewDecoder(f).Decode(&tasks)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}

	id := len(tasks) + 1
	now := time.Now().Format("2006-01-02 03-04 PM")

	var task = Task{Title: title, Status: status, ID: id, CreatedAt: now, UpdatedAt: now}

	tasks = append(tasks, task)
	newData, _ := json.MarshalIndent(tasks, "", " ")
	os.WriteFile("tracker.json", newData, 0644)

	fmt.Println("New task successfully added.")
}

// delete task
func delete_task(num int) int {
	f, err := os.OpenFile(
		"tracker.json",
		os.O_CREATE|os.O_RDONLY,
		0644,
	)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return -1
	}
	defer f.Close()

	var tasks []Task
	err = json.NewDecoder(f).Decode(&tasks)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return -1
	}

	if len(tasks) == 0 {
		return 2
	}

	if num < 1 || num > len(tasks) {
		return 0
	}

	if num == len(tasks) {
		tasks = tasks[:num-1]
	} else {
		tasks = append(tasks[:num-1], tasks[num:]...)
	}

	for i := range tasks {
		if tasks[i].ID > num {
			tasks[i].ID -= 1
		}
	}

	newData, _ := json.MarshalIndent(tasks, "", " ")
	os.WriteFile("tracker.json", newData, 0644)

	return 1
}

// change task title
func modify_task(taskNum int, newTitle string) int {
	f, err := os.OpenFile(
		"tracker.json",
		os.O_CREATE|os.O_RDONLY,
		0644,
	)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	defer f.Close()

	var tasks []Task
	err = json.NewDecoder(f).Decode(&tasks)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return -1
	}

	if taskNum < 1 || taskNum > len(tasks) {
		return 0
	}

	tasks[taskNum-1].Title = newTitle
	now := time.Now().Format("2006-01-02 03-04 PM")
	tasks[taskNum-1].UpdatedAt = now

	newData, _ := json.MarshalIndent(tasks, "", " ")
	os.WriteFile("tracker.json", newData, 0644)
	return 1
}

// change task status
func update_task(taskNum int, statusNum int) int {
	f, err := os.OpenFile(
		"tracker.json",
		os.O_CREATE|os.O_RDONLY,
		0644,
	)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer f.Close()

	var tasks []Task
	err = json.NewDecoder(f).Decode(&tasks)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return -1
	}

	switch statusNum {
	case 1:
		tasks[taskNum-1].Status = "todo"
	case 2:
		tasks[taskNum-1].Status = "on-progress"
	case 3:
		tasks[taskNum-1].Status = "done"
	default:
		return 0
	}

	now := time.Now().Format("2006-01-02 03-04 PM")
	tasks[taskNum-1].UpdatedAt = now

	newData, _ := json.MarshalIndent(tasks, "", " ")
	os.WriteFile("tracker.json", newData, 0644)
	return 1
}

// print list of tasks
func list_tasks() {
	var tasks []Task

	f, err := os.OpenFile(
		"tracker.json",
		os.O_CREATE|os.O_RDONLY,
		0644,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&tasks)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("(There is no task!)")
		return
	}

	for _, element := range tasks {
		fmt.Printf("%d. %s - %s\ncreate time: %s - last modify: %s\n\n", element.ID, element.Title, element.Status, element.CreatedAt, element.UpdatedAt)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var usrChoice int
	for ok := true; ok; ok = (usrChoice != 6) {
		CallClear()
		menu()
		scanner.Scan()
		choice := scanner.Text()
		usrChoice, _ = strconv.Atoi(strings.TrimSpace(choice))

		switch usrChoice {
		case 1:
			fmt.Printf("enter your task: ")
			scanner.Scan()
			desc := scanner.Text()
			fmt.Printf("enter the task status(1. todo 2. in-progress 3. done): ")
			scanner.Scan()
			statusInt, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))

			switch statusInt {
			case 1:
				add_task(desc, "todo")
			case 2:
				add_task(desc, "in-progress")
			case 3:
				add_task(desc, "done")
			default:
				fmt.Println("Invalid response. Please try again. Press Enter to continue...")
				scanner.Scan()
			}

		case 2:
			list_tasks()
			fmt.Println("Which task you want to modify(enter the number): ")
			scanner.Scan()
			taskNum, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			fmt.Println("What you want to change task to: ")
			scanner.Scan()
			newTitle := scanner.Text()
			err := modify_task(taskNum, newTitle)
			switch err {
			case 1:
				fmt.Println("Task successfully changed. Press enter to continue...")
				scanner.Scan()
			case 0:
				fmt.Println("Something went wrong, try again. Press enter to continue...")
				scanner.Scan()
			case -1:
				fmt.Println("Error occurred. Try again. Press enter to continue...")
				scanner.Scan()
			}

		case 3:
			list_tasks()
			fmt.Println("Which one you want to delete(enter the number): ")
			scanner.Scan()
			choice := scanner.Text()
			taskId, _ := strconv.Atoi(strings.TrimSpace(choice))
			err := delete_task(taskId)
			switch err {
			case 1:
				fmt.Println("Task deleted successfully. Press enter to continue...")
				scanner.Scan()
			case 0:
				fmt.Println("Invalid input, try again. Press enter to continue...")
				scanner.Scan()
			case -1:
				fmt.Println("Error occurred, try again. Press enter to continue...")
				scanner.Scan()
			case 2:
				fmt.Println("File is empty. Try adding new tasks. Press enter to continue...")
				scanner.Scan()
			}

		case 4:
			list_tasks()
			fmt.Println("Which task you want to update(enter the number): ")
			scanner.Scan()
			taskNum, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			fmt.Println("Chnage to (1. todo 2. on-progress 3. done): ")
			scanner.Scan()
			statusNum, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			err := update_task(taskNum, statusNum)
			switch err {
			case 1:
				fmt.Println("Task status changed successfully. Press enter to continue...")
				scanner.Scan()
			case 0:
				fmt.Println("Something went wrong, try again. Press enter to continue...")
				scanner.Scan()
			case -1:
				fmt.Println("Error occurred. Try again. Press enter to continue...")
				scanner.Scan()
			}

		case 5:
			list_tasks()
			fmt.Println("Press enter to continue...")
			scanner.Scan()
		}
	}
}
