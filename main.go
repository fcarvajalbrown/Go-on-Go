package main

import (
	"go-game/game"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// In-memory storage for games (use database in production)
var games = make(map[string]*game.Board)

func main() {
	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Serve static files (HTML, CSS, JS for game board)
	e.Static("/", "static")

	// WebSocket endpoint for real-time game moves
	e.GET("/ws", handleWebSocket)

	// REST API endpoints
	e.POST("/game/new", newGame)       // Create new game
	e.GET("/game/:id", getGame)        // Get game state
	e.POST("/game/:id/move", makeMove) // Make a move

	// Start server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}

// WebSocket handler for real-time communication
func handleWebSocket(c echo.Context) error {
	return nil
}

// Create new Go game
func newGame(c echo.Context) error {
	// Create a new 19x19 Go board
	board := game.NewBoard(19)

	// Store it with a fixed ID for now (use UUID in production)
	gameID := "local"
	games[gameID] = board

	// Return the board state
	return c.JSON(http.StatusOK, board)
}

// Get current game state
func getGame(c echo.Context) error {
	gameID := c.Param("id")

	// Find the game
	board, exists := games[gameID]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Game not found"})
	}

	return c.JSON(http.StatusOK, board)
}

// Move request structure
type MoveRequest struct {
	Position int  `json:"position"` // Board position (0-360 for 19x19)
	Pass     bool `json:"pass"`     // True if player wants to pass
}

// Process player move
func makeMove(c echo.Context) error {
	gameID := c.Param("id")

	// Find the game
	board, exists := games[gameID]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Game not found"})
	}

	// Parse the move request
	var moveReq MoveRequest
	if err := c.Bind(&moveReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	// Handle pass move
	if moveReq.Pass {
		board.Pass()
		return c.JSON(http.StatusOK, board)
	}

	// Validate position range
	if moveReq.Position < 0 || moveReq.Position >= board.Size*board.Size {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Position out of bounds"})
	}

	// Attempt to make the move
	if err := board.MakeMove(moveReq.Position); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Return updated board state
	return c.JSON(http.StatusOK, board)
}
