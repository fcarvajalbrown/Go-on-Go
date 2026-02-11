package game

import "fmt"

// Board represents the game state of a Go board
// Go is played on a 19x19 grid with complex rules for capturing and scoring
type Board struct {
	// Size of the board (typically 19, but can be 9, 13, or other sizes)
	Size int

	// Grid stores the current state of each intersection
	// 0 = empty, 1 = black stone, 2 = white stone
	// We use a 1D slice for efficiency: position = row*size + col
	Grid []int

	// CurrentPlayer tracks whose turn it is (1 = black, 2 = white)
	// Black always plays first in Go
	CurrentPlayer int

	// CapturedStones tracks how many stones each player has captured
	// Index 0 is unused, index 1 = stones captured by black, index 2 = stones captured by white
	CapturedStones [3]int

	// Ko represents the "Ko rule" - prevents infinite loops
	// Stores the board position from the previous move to prevent immediate recapture
	Ko []int

	// MoveHistory stores all moves made in the game for game review and undo functionality
	MoveHistory []Move
}

// Move represents a single move in the game
type Move struct {
	// Player who made the move (1 = black, 2 = white)
	Player int

	// Position on the board (row*size + col), -1 for pass
	Position int

	// CapturedPositions stores which stones were captured by this move
	// Needed for proper undo functionality and Ko rule enforcement
	CapturedPositions []int
}

// NewBoard creates a new Go board with the specified size
// Standard sizes are 9x9 (beginner), 13x13 (intermediate), 19x19 (professional)
func NewBoard(size int) *Board {
	return &Board{
		Size:           size,
		Grid:           make([]int, size*size), // All positions start empty (0)
		CurrentPlayer:  1,                      // Black plays first
		CapturedStones: [3]int{0, 0, 0},        // No captured stones initially
		Ko:             nil,                    // No Ko situation initially
		MoveHistory:    make([]Move, 0),        // Empty move history
	}
}

// IsValidPosition checks if a coordinate is within the board boundaries
func (b *Board) IsValidPosition(row, col int) bool {
	return row >= 0 && row < b.Size && col >= 0 && col < b.Size
}

// GetPosition converts row, col coordinates to 1D array index
func (b *Board) GetPosition(row, col int) int {
	return row*b.Size + col
}

// GetCoordinates converts 1D array index back to row, col coordinates
func (b *Board) GetCoordinates(position int) (int, int) {
	return position / b.Size, position % b.Size
}

// IsEmpty checks if a position on the board is empty
func (b *Board) IsEmpty(position int) bool {
	return b.Grid[position] == 0
}

// GetStone returns the color of stone at a position (0=empty, 1=black, 2=white)
func (b *Board) GetStone(position int) int {
	return b.Grid[position]
}

// GetNeighbors returns all adjacent positions (up, down, left, right)
// In Go, only orthogonally adjacent intersections matter for captures and groups
func (b *Board) GetNeighbors(position int) []int {
	row, col := b.GetCoordinates(position)
	neighbors := make([]int, 0, 4) // Maximum 4 neighbors

	// Check all four directions: up, down, left, right
	directions := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, dir := range directions {
		newRow, newCol := row+dir[0], col+dir[1]
		if b.IsValidPosition(newRow, newCol) {
			neighbors = append(neighbors, b.GetPosition(newRow, newCol))
		}
	}

	return neighbors
}

// GetGroup finds all stones connected to a given position
// In Go, stones of the same color that are orthogonally connected form a "group"
// Groups live or die together based on their collective liberties
func (b *Board) GetGroup(position int) []int {
	color := b.GetStone(position)
	if color == 0 {
		return nil // Empty position has no group
	}

	visited := make(map[int]bool)
	group := make([]int, 0)

	// Use depth-first search to find all connected stones of the same color
	var dfs func(int)
	dfs = func(pos int) {
		if visited[pos] || b.GetStone(pos) != color {
			return
		}

		visited[pos] = true
		group = append(group, pos)

		// Recursively check all neighbors
		for _, neighbor := range b.GetNeighbors(pos) {
			dfs(neighbor)
		}
	}

	dfs(position)
	return group
}

// GetLiberties counts the number of empty adjacent positions for a group
// "Liberties" are empty intersections adjacent to a group
// A group with no liberties is captured and removed from the board
func (b *Board) GetLiberties(group []int) int {
	liberties := make(map[int]bool) // Use map to avoid counting same liberty twice

	for _, pos := range group {
		for _, neighbor := range b.GetNeighbors(pos) {
			if b.IsEmpty(neighbor) {
				liberties[neighbor] = true
			}
		}
	}

	return len(liberties)
}

// WouldBeSuicide checks if placing a stone would be suicide
// Suicide is placing a stone that would immediately have no liberties
// This is illegal unless the move captures opponent stones
func (b *Board) WouldBeSuicide(position int, player int) bool {
	// Temporarily place the stone
	originalStone := b.Grid[position]
	b.Grid[position] = player

	// Check if this creates a group with liberties
	group := b.GetGroup(position)
	liberties := b.GetLiberties(group)

	// Restore original state
	b.Grid[position] = originalStone

	// If the group would have no liberties, it's potentially suicide
	if liberties > 0 {
		return false
	}

	// Check if this move would capture opponent stones
	// If it captures opponent stones, it's not suicide even with no liberties
	opponent := 3 - player // Convert 1->2, 2->1
	for _, neighbor := range b.GetNeighbors(position) {
		if b.GetStone(neighbor) == opponent {
			opponentGroup := b.GetGroup(neighbor)
			if b.GetLiberties(opponentGroup) == 1 {
				return false // This move would capture, so not suicide
			}
		}
	}

	return true // No liberties and no captures = suicide
}

// IsValidMove checks if a move is legal according to Go rules
func (b *Board) IsValidMove(position int) bool {
	// Move must be on an empty intersection
	if !b.IsEmpty(position) {
		return false
	}

	// Move cannot be suicide
	if b.WouldBeSuicide(position, b.CurrentPlayer) {
		return false
	}

	// Move cannot violate Ko rule (immediate recapture)
	if b.Ko != nil && len(b.Ko) == len(b.Grid) {
		// Temporarily make the move and check if it recreates the Ko position
		b.Grid[position] = b.CurrentPlayer
		captures := b.processCaptures(position)

		isKo := true
		for i, stone := range b.Grid {
			if stone != b.Ko[i] {
				isKo = false
				break
			}
		}

		// Restore board state
		b.Grid[position] = 0
		for _, capturedPos := range captures {
			b.Grid[capturedPos] = 3 - b.CurrentPlayer
		}

		if isKo {
			return false // Ko rule violation
		}
	}

	return true
}

// processCaptures handles capturing opponent groups that have no liberties
// Returns the positions of captured stones
func (b *Board) processCaptures(position int) []int {
	opponent := 3 - b.CurrentPlayer // Convert 1->2, 2->1
	captured := make([]int, 0)

	// Check all adjacent opponent groups
	for _, neighbor := range b.GetNeighbors(position) {
		if b.GetStone(neighbor) == opponent {
			group := b.GetGroup(neighbor)
			if b.GetLiberties(group) == 0 {
				// This group has no liberties, capture it
				for _, pos := range group {
					b.Grid[pos] = 0 // Remove stone
					captured = append(captured, pos)
				}
				b.CapturedStones[b.CurrentPlayer] += len(group)
			}
		}
	}

	return captured
}

// MakeMove places a stone on the board and handles all game logic
func (b *Board) MakeMove(position int) error {
	if !b.IsValidMove(position) {
		return fmt.Errorf("invalid move at position %d", position)
	}

	// Save current board state for Ko rule
	previousBoard := make([]int, len(b.Grid))
	copy(previousBoard, b.Grid)

	// Place the stone
	b.Grid[position] = b.CurrentPlayer

	// Process captures
	captured := b.processCaptures(position)

	// Record the move
	move := Move{
		Player:            b.CurrentPlayer,
		Position:          position,
		CapturedPositions: captured,
	}
	b.MoveHistory = append(b.MoveHistory, move)

	// Update Ko position
	b.Ko = previousBoard

	// Switch players
	b.CurrentPlayer = 3 - b.CurrentPlayer

	return nil
}

// Pass allows a player to skip their turn
func (b *Board) Pass() {
	move := Move{
		Player:   b.CurrentPlayer,
		Position: -1, // -1 indicates a pass
	}
	b.MoveHistory = append(b.MoveHistory, move)

	// Switch players
	b.CurrentPlayer = 3 - b.CurrentPlayer
}

// IsGameOver checks if the game has ended (both players passed consecutively)
func (b *Board) IsGameOver() bool {
	if len(b.MoveHistory) < 2 {
		return false
	}

	// Game ends when both players pass in succession
	lastMove := b.MoveHistory[len(b.MoveHistory)-1]
	secondLastMove := b.MoveHistory[len(b.MoveHistory)-2]

	return lastMove.Position == -1 && secondLastMove.Position == -1
}
