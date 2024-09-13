package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Maze struct to represent the maze map and tractor
type Maze struct {
	mapData       [][]string
	tractor       Tractor
	columnsNumber int
	rowsNumber    int
}

// Tractor struct to keep track of the tractor's position
type Tractor struct {
	x, y int
}

type User struct {
	x              int
	y              int
	victoryCounter int
}

// Helper function to check if a number is even
func isEven(n int) bool {
	return n%2 == 0
}

// Helper function to get a random element from a slice of integers
func getRandomFromInt(slice []int) int {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	index := rand.Intn(len(slice))   // Generate a random index
	return slice[index]              // Return the element at the random index
}

func getRandomFromString(slice []string) string {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	index := rand.Intn(len(slice))   // Generate a random index
	return slice[index]              // Return the element at the random index
}

// Helper function to initialize the maze
func generateMaze(columnsNumber, rowsNumber int) [][]string {
	maze := Maze{
		mapData:       make([][]string, rowsNumber),
		columnsNumber: columnsNumber,
		rowsNumber:    rowsNumber,
	}

	// Fill the maze with walls
	for y := 0; y < rowsNumber; y++ {
		row := make([]string, columnsNumber)
		for x := 0; x < columnsNumber; x++ {
			row[x] = "▉"
		}
		maze.mapData[y] = row
	}

	// Choose random even coordinates for the tractor
	evenXCoords := []int{}
	evenYCoords := []int{}
	for i := 0; i < columnsNumber; i++ {
		if isEven(i) {
			evenXCoords = append(evenXCoords, i)
		}
	}
	for i := 0; i < rowsNumber; i++ {
		if isEven(i) {
			evenYCoords = append(evenYCoords, i)
		}
	}
	maze.tractor = Tractor{
		x: getRandomFromInt(evenXCoords),
		y: getRandomFromInt(evenYCoords),
	}

	// Clear the cell where the tractor is initially placed
	maze.setField(maze.tractor.x, maze.tractor.y, " ")

	// Generate the maze
	for !maze.isMaze() {
		maze.moveTractor()
	}

	return maze.mapData
}

// Helper function to get the value of a cell in the maze
func (m *Maze) getField(x, y int) string {
	if x < 0 || x >= m.columnsNumber || y < 0 || y >= m.rowsNumber {
		return ""
	}
	return m.mapData[y][x]
}

// Helper function to set the value of a cell in the maze
func (m *Maze) setField(x, y int, value string) {
	if x >= 0 && x < m.columnsNumber && y >= 0 && y < m.rowsNumber {
		m.mapData[y][x] = value
	}
}

// Function to move the tractor in the maze
func (m *Maze) moveTractor() {
	directs := []string{}

	if m.tractor.x > 0 {
		directs = append(directs, "left")
	}
	if m.tractor.x < m.columnsNumber-2 {
		directs = append(directs, "right")
	}
	if m.tractor.y > 0 {
		directs = append(directs, "up")
	}
	if m.tractor.y < m.rowsNumber-2 {
		directs = append(directs, "down")
	}

	direct := getRandomFromString(directs)

	switch direct {
	case "left":
		if m.getField(m.tractor.x-2, m.tractor.y) == "▉" {
			m.setField(m.tractor.x-1, m.tractor.y, " ")
			m.setField(m.tractor.x-2, m.tractor.y, " ")
		}
		m.tractor.x -= 2
	case "right":
		if m.getField(m.tractor.x+2, m.tractor.y) == "▉" {
			m.setField(m.tractor.x+1, m.tractor.y, " ")
			m.setField(m.tractor.x+2, m.tractor.y, " ")
		}
		m.tractor.x += 2
	case "up":
		if m.getField(m.tractor.x, m.tractor.y-2) == "▉" {
			m.setField(m.tractor.x, m.tractor.y-1, " ")
			m.setField(m.tractor.x, m.tractor.y-2, " ")
		}
		m.tractor.y -= 2
	case "down":
		if m.getField(m.tractor.x, m.tractor.y+2) == "▉" {
			m.setField(m.tractor.x, m.tractor.y+1, " ")
			m.setField(m.tractor.x, m.tractor.y+2, " ")
		}
		m.tractor.y += 2
	}
}

// Function to check if the maze is complete
func (m *Maze) isMaze() bool {
	for x := 0; x < m.columnsNumber; x++ {
		for y := 0; y < m.rowsNumber; y++ {
			if isEven(x) && isEven(y) && m.getField(x, y) == "▉" {
				return false
			}
		}
	}

	return true
}

func moveForward(c *gin.Context) {
	if user.x == rowsNumber-1 && user.y != columnsNumber-2 {
		// нижняя граница лабиринта - движение невозможно
		c.JSON(http.StatusConflict, gin.H{"message": "The lower limit has been reached"})
		return
	}

	if maze[user.x+1][user.y] == "▉" {
		c.JSON(http.StatusConflict, gin.H{"message": "The wall is ahead, you can't move"})
		return
	}

	user.x += 1
	if user.x == rowsNumber-1 && user.y == columnsNumber-2 {
		user.victoryCounter += 1
		c.JSON(http.StatusOK, gin.H{"message": "Your " + strconv.Itoa(user.victoryCounter) + " win! Game is restarted"})
		startGame()
		return
	}

	if user.x == rowsNumber-1 && maze[rowsNumber-1][columnsNumber-1] == " " && user.y == columnsNumber-1 {
		user.victoryCounter += 1
		c.JSON(http.StatusOK, gin.H{"message": "Your " + strconv.Itoa(user.victoryCounter) + " win! Game is restarted"})
		startGame()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "You moved forward, your current position - (" + strconv.Itoa(user.x) +
		";" + strconv.Itoa(user.y) + ")"})
	return
}

func moveBackward(c *gin.Context) {

	if user.x == 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "You are at the first row, you can't move back"})
		return
	}

	if maze[user.x-1][user.y] == "▉" {
		c.JSON(http.StatusConflict, gin.H{"message": "The wall is behind, you can't move"})
		return
	}

	user.x -= 1
	c.JSON(http.StatusOK, gin.H{"message": "You moved backward, your current position - (" + strconv.Itoa(user.x) +
		";" + strconv.Itoa(user.y) + ")"})
	return
}

func moveLeft(c *gin.Context) {

	if user.y == 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "You are at the left edge, you can't move left"})
		return
	}

	if maze[user.x][user.y-1] == "▉" {
		c.JSON(http.StatusConflict, gin.H{"message": "The wall is on the left, you can't move"})
		return
	}

	user.y -= 1
	c.JSON(http.StatusOK, gin.H{"message": "You moved left, your current position - (" + strconv.Itoa(user.x) +
		";" + strconv.Itoa(user.y) + ")"})
	return

}

func moveRight(c *gin.Context) {

	if user.y == columnsNumber-1 {
		c.JSON(http.StatusConflict, gin.H{"message": "You are at the right edge, you can't move right"})
		return
	}

	if maze[user.x][user.y+1] == "▉" {
		c.JSON(http.StatusConflict, gin.H{"message": "The wall is on the right, you can't move"})
		return
	}

	user.y += 1
	if user.x == rowsNumber-1 && user.y == columnsNumber-2 {
		user.victoryCounter += 1
		c.JSON(http.StatusOK, gin.H{"message": "Your " + strconv.Itoa(user.victoryCounter) + " win! Game is restarted"})
		startGame()
		return
	}

	if user.x == rowsNumber-1 && maze[rowsNumber-1][columnsNumber-1] == " " && user.y == columnsNumber-1 {
		user.victoryCounter += 1
		c.JSON(http.StatusOK, gin.H{"message": "Your " + strconv.Itoa(user.victoryCounter) + " win! Game is restarted"})
		startGame()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "You moved right, your current position - (" + strconv.Itoa(user.x) +
		";" + strconv.Itoa(user.y) + ")"})
	return
}

func lookAround(c *gin.Context) {

	left := ""
	right := ""
	back := ""
	forward := ""

	if user.y == 0 {
		left = "left edge of labyrinth"
	} else {
		left = maze[user.x][user.y-1]
		if left == " " {
			left = "empty"
		} else {
			left = "wall"
		}
	}

	if user.y == columnsNumber-1 {
		right = "right edge of labyrinth"
	} else {
		right = maze[user.x][user.y+1]
		if right == " " {
			right = "empty"
		} else {
			right = "wall"
		}
	}

	if user.x == 0 {
		back = "start edge of labyrinth"
	} else {
		back = maze[user.x-1][user.y]
		if back == " " {
			back = "empty"
		} else {
			back = "wall"
		}
	}

	if user.x == rowsNumber-1 {
		forward = "end edge of labyrinth"
	} else {
		forward = maze[user.x+1][user.y]
		if forward == " " {
			forward = "empty"
		} else {
			forward = "wall"
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "On the left - " + left + ". On the right - " +
		right + ". On the back - " + back + ". On the forward - " + forward})
}

func startGame() {
	user.x, user.y = 0, 0
	maze = generateMaze(rowsNumber, columnsNumber)
	maze[rowsNumber-1][columnsNumber-2] = " "
	// Output the final maze
	fmt.Println("Generated maze:")
	for _, row := range maze {
		fmt.Println(row)
	}
}

var columnsNumber = 10
var rowsNumber = 10
var maze = make([][]string, rowsNumber)
var user = User{0, 0, 0}

func main() {
	startGame()
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to labyrinth!"})
	})
	router.GET("/forward", moveForward)
	router.GET("/backward", moveBackward)
	router.GET("/left", moveLeft)
	router.GET("/right", moveRight)
	router.GET("/look", lookAround)
	router.GET("/victory", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "You won " + strconv.Itoa(user.victoryCounter) + " times"})
	})

	router.Run("localhost:8080")
}
