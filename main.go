package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Task struct
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdat"`
	UpdatedAt   string `json:"updatedat"`
}

const dataFile = "tasks.json"

//-----------Utility functions-------------------
/*loadTasks()       ← reads tasks.json into memory
saveTasks()       ← writes tasks back to tasks.json
generateID()      ← generates unique task ID
findTaskByID()    ← finds task by ID in array
isValidStatus()   ← checks if status is valid (todo/in-progress/done)*/

// 1. loadTasks-load tasks from json file
func loadTasks() ([]Task, error) {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		os.WriteFile(dataFile, []byte("[]"), 0644)
	}
	file, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	json.Unmarshal(file, &tasks)
	return tasks, nil
}

// 2.saveTasks-save tasks to json file
func saveTasks(tasks []Task) error {
	data, _ := json.MarshalIndent(tasks, "", " ")
	return os.WriteFile(dataFile, data, 0644)
}

// 3.generateID-generate a new unique task id
func generateID(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}
	return tasks[len(tasks)-1].ID + 1
}

// 4.isValidStatus-Validate status input
func isValidStatus(status string) bool {
	return status == "todo" || status == "in-progress" || status == "done"

}

// 5.findTaskByID - find task by id
func findTaskByID(tasks []Task, id int) *Task {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i]
		}
	}
	return nil
}

// Comment funcs
/*1.add Task
2.list Task
3.update Task
4.delete Task
*/
//1.addTask
func addTask(tasks []Task, description string) {
	newID := generateID(tasks)
	now := time.Now().Format("2006-01-02T15:05:06")
	newTask := Task{
		ID:          newID,
		Description: description,
		Status:      "todo",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	tasks = append(tasks, newTask)
	saveTasks(tasks)
	fmt.Printf("Task added successfully(ID: %d)\n", newID)
}

// 2.ListTasks (with optional status filter)
func listTasks(tasks []Task, filter string) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}
	for _, t := range tasks {
		if filter == "" || t.Status == filter {
			fmt.Printf("%d.[%s]%s(Created:%s,Updated:%s)\n",
				t.ID, t.Status, t.Description, t.CreatedAt, t.UpdatedAt)
		}
	}
}

// 3.UpdateTask -  Update task status
func updateTask(tasks []Task, id int, newStatus string) {
	task := findTaskByID(tasks, id)
	if task == nil {
		fmt.Println("Task not found")
		return
	}
	if !isValidStatus(newStatus) {
		fmt.Println("Invalid status.Use:todo,in-progress,done")
		return
	}
	task.Status = newStatus
	task.UpdatedAt = time.Now().Format("2006-04-08T12:05:03")
	saveTasks(tasks)
	fmt.Printf("Task %d updated to %s \n", id, newStatus)
}

// 4.Delete Task- Delete a Task
func deleteTask(tasks []Task, id int) {
	newTasks := []Task{}
	deleted := false
	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		} else {
			deleted = true
		}
	}
	if deleted {
		saveTasks(newTasks)
		fmt.Printf("Task %d deleted\n", id)
	} else {
		fmt.Printf("Task not found")
	}
}

// Main function
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-cli [add|list|update|delete] arguments")
		return
	}
	command := os.Args[1]
	tasks, _ := loadTasks()

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli add \"Task description\"")
			return
		}
		addTask(tasks, os.Args[2])

	case "list":
		filter := ""
		if len(os.Args) >= 3 {
			filter = os.Args[2]

		}
		listTasks(tasks, filter)

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: task-cli update <id> <status>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID")
			return
		}
		updateTask(tasks, id, os.Args[3])

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli delete <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID")
			return
		}
		deleteTask(tasks, id)

	default:
		fmt.Println("Usage:")
		fmt.Println("  add \"task description\"   → Add a new task")
		fmt.Println("  update <id> <status>      → Update task status (todo/in-progress/done)")
		fmt.Println("  delete <id>               → Delete a task")
		fmt.Println("  list                      → List all tasks")
		fmt.Println("  list <status>             → List tasks by status")

	}
}
