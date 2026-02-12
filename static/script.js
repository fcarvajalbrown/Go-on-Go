let gameState = null;
let gameId = 'local'; // For now, use a fixed game ID

// Initialize the board UI
function initBoard() {
    const board = document.getElementById('board');
    board.innerHTML = ''; // Clear existing board
    
    // Star points positions for 19x19 board (traditional Go board markers)
    const starPoints = [60, 66, 72, 174, 180, 186, 288, 294, 300]; // 3-3, 3-9, 3-15, 9-3, 9-9, 9-15, 15-3, 15-9, 15-15
    
    for(let i = 0; i < 361; i++) {
        const cell = document.createElement('div');
        cell.className = 'cell';
        
        // Add star points
        if (starPoints.includes(i)) {
            cell.classList.add('star-point');
        }
        
        cell.dataset.position = i;
        cell.onclick = () => makeMove(i);
        board.appendChild(cell);
    }
}

// Create a new game
async function newGame() {
    try {
        const response = await fetch('/game/new', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        });
        
        if (response.ok) {
            gameState = await response.json();
            updateUI();
            setStatus('New game started!');
        } else {
            setStatus('Failed to create new game');
        }
    } catch (error) {
        setStatus('Error: ' + error.message);
    }
}

// Make a move at the specified position
async function makeMove(position) {
    if (!gameState) {
        setStatus('Please start a new game first');
        return;
    }
    
    try {
        const response = await fetch(`/game/${gameId}/move`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ position: position })
        });
        
        if (response.ok) {
            gameState = await response.json();
            updateUI();
            setStatus(`Move played at position ${position}`);
        } else {
            const error = await response.text();
            setStatus('Invalid move: ' + error);
        }
    } catch (error) {
        setStatus('Error: ' + error.message);
    }
}

// Pass turn
async function pass() {
    if (!gameState) {
        setStatus('Please start a new game first');
        return;
    }
    
    try {
        const response = await fetch(`/game/${gameId}/move`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ pass: true })
        });
        
        if (response.ok) {
            gameState = await response.json();
            updateUI();
            setStatus('Player passed');
        } else {
            setStatus('Failed to pass');
        }
    } catch (error) {
        setStatus('Error: ' + error.message);
    }
}

// Update the UI with current game state
function updateUI() {
    if (!gameState) return;
    
    // Update board stones
    const cells = document.querySelectorAll('.cell');
    gameState.Grid.forEach((stone, index) => {
        const cell = cells[index];
        cell.className = 'cell';
        if (starPoints.includes(index)) {
            cell.classList.add('star-point');
        }
        if (stone === 1) cell.classList.add('black');
        else if (stone === 2) cell.classList.add('white');
    });
    
    // Update current player
    const currentPlayerEl = document.getElementById('currentPlayer');
    const isBlack = gameState.CurrentPlayer === 1;
    currentPlayerEl.textContent = isBlack ? 'Black' : 'White';
    currentPlayerEl.className = isBlack ? 'current-player current-black' : 'current-player current-white';
    
    // Update captured stones
    document.getElementById('capturedBlack').textContent = gameState.CapturedStones[1] || 0;
    document.getElementById('capturedWhite').textContent = gameState.CapturedStones[2] || 0;
    
    // Check for game over
    if (gameState.IsGameOver) {
        setStatus('Game Over! Both players passed.');
    }
}

// Display status messages
function setStatus(message) {
    document.getElementById('status').textContent = message;
}

// Star points for reference (needed in updateUI)
const starPoints = [60, 66, 72, 174, 180, 186, 288, 294, 300];

// Initialize the game on page load
document.addEventListener('DOMContentLoaded', function() {
    initBoard();
    setStatus('Click "New Game" to start playing');
});