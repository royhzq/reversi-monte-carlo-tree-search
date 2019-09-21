// Reversi/Othello Game engine 
// 2019 Roy Hung

package main

import (
	"fmt"
	"time"
	"math"
	"strconv"
	"math/rand"
)

type Position struct {
	i int // ith row in a reversi/othello board
	j int // jth column in a reversi/othello board 
}

type strPosition struct {
	i string // ith row in a reversi/othello board (string)
	j string // jth column in a reversi/othello board (string)
}

// Eight possible directions to traverse on the board
var Directions = []Position{
    Position{-1,0},
    Position{-1,1},
    Position{0,1},
    Position{1,1},
    Position{1,0},
    Position{1,-1},
    Position{0,-1},
    Position{-1,-1},
}

// Corner pieces in reversi/othello are considered high in value
// For heuristics to supplement UCT to select child node
var corners = []Position{
	{0,0},{0,7},{7,0},{7,7},
}

// Generally considered bad positions 
// These are positions adjacent to the corners
// For heuristics to supplement UCT to select child node
var badPositions = []Position{
	{0,1},{1,0},
	{6,0},{7,1},
	{7,6},{6,7},
	{0,6},{1,7},
}

// Generally considered very bad positions 
// These are positions that give corners away
// For heuristics to supplement UCT to select child node
var veryBadPositions = []Position{
	{1,1},{1,6},{6,1},{6,6},
}

func (position Position) PrintPrettifyNotation() strPosition {
	// Converts Position from (row, column) notation 
	// to standard Othello notation 
	// Example: 
	//     (0,0) -> "A1", Top left corner
	//     (0,7) -> "H1", Top right corner
	alphabet := []string{"A","B","C","D","E","F","G","H"}
	return strPosition{alphabet[position.j], strconv.Itoa(position.i+1)}
}

func posInSlice(a Position, list []Position) bool {
	// Check if a Position is in a Slice
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

// The Reversi/Othello board
type Board struct {
	length int            // Max length of board (i.e., 8 for standard board size)
	board [8][8]int       // State of board
	filled []Position     // Slice of all spaces filled up by a piece
	empty []Position      // Slice of all spaces that remain empty
	neighbours []Position // Slice of all empty neighbours of pieces
	validSpace []Position // Slice of all valid moves for the player turn
	blackScore int        // Total number of Black pieces on board(1)  
	whiteScore int        // Total number of White pieces on board (-1) 
	winner int            // Winner of game - Black (1), White (-1), Draw (99), Undetermined (0). Undetermined is default
	turn int              // Whose turn is it (1 for Black, -1 for White)
}

func (X Board) Show() {
	// Display board in terminal
	dim := len(X.board) -1
	for i := 0; i <= dim; i++ {
		for j:= 0; j <= dim; j++ {
			showPiece := "   "
			if  X.board[i][j] == 1 {
				showPiece = " X "
			}
			if X.board[i][j] == -1 {
				showPiece = " O "
			}
			fmt.Printf("%s", showPiece )
		}
		fmt.Println("\n")
	}
}

func (X Board) inRange(space Position) bool {
	// Check if Position can exist in the board 
	// Example: 
	//     Position{12, 8} is out of range and invalid in an 8x8 board
	// Return false if Position is out of range
	if space.i < X.length && 
	   space.j < X.length && 
	   space.i >= 0 && 
	   space.j>=0 {
	   	return true
	} else {
		return false
	}
} 

func (X *Board) getNeighbour(piece Position) []Position {
	// Obtain a Slice of immediate neighbouring empty spaces for a given Position	
	neighbours := []Position{}
	for _, dir := range Directions {

		// Iterate through all surrounding spaces for given piece
		currPos := Position{piece.i + dir.i, piece.j + dir.j}
		if X.inRange(currPos) && X.board[currPos.i][currPos.j] == 0 {
			neighbours = append(neighbours, currPos)
		}
	}
	return neighbours
}

func (X *Board) initNeighbours() {
	// Obtain all neighbours on board
	filled := X.filled 
	neighbours := []Position{}     // List to store all neighbours on board
	pieceNeighbour := []Position{} // List to store neighbours of the current piece

	// Loop through all filled spaces and get their neighbours
	// Add neighbours to neighbours list if position does not already exist in list
	for _, piece := range filled {
		pieceNeighbour = X.getNeighbour(piece)
		for _, n := range pieceNeighbour {
			if posInSlice(n, neighbours) == false {
				neighbours = append(neighbours, n)
			}
		}
	}
	X.neighbours = neighbours
}

func (X *Board) getScores() (int, int) {
	// Count all black/white pieces
	// To calculates overall score for (black, white)
	blackScore := 0
	whiteScore := 0
	for i := 0; i < X.length; i++ { 
		for j := 0; j < X.length; j++ {
			if X.board[i][j] == 1 {
				blackScore += 1
			}
			if X.board[i][j] == -1 {
				whiteScore += 1
			}
		}
	}
	return blackScore, whiteScore
}

func (X *Board) Setup() {
	// Initialize parameters of a Board
	// Called each time a new board is created
	X.blackScore = 0
	X.whiteScore = 0
	X.filled = []Position{}
	X.empty = []Position{}
	X.neighbours = []Position{}
	X.validSpace = []Position{}

	// Initialize empty and filled parameters
	for i := 0; i < X.length; i++ { 
		for j := 0; j < X.length; j++ {
			if X.board[i][j] == 0 {
				X.empty = append(X.empty, Position{i,j})
			} else {
				X.filled = append(X.filled, Position{i,j})
			}
		}
	}

	// Initialize set of all neighbours 
	X.initNeighbours()

	// Initialize validSpace after neighbours initialized
	X.validSpace = X.getAllValid()

	// Initialize starting scores 
	X.blackScore, X.whiteScore = X.getScores()

	// Initialize winner as Undetermined - 0
	X.winner = 0
}

func (X *Board) checkValidDir(iDir int, jDir int, space Position) bool {
	// Check, in the direction (iDir, jDir), whether pieces will flip
	// Direction is defined by (iDir, jDir), for example,
	// Example:
	//    iDir=1, jDir=1, we are checking the N-E direction of the given Position (space)
	//    iDir=0, jDir=-1, we are checking the West direction of the given Position (space)
	// First Check if space is empty. space has to be empty (0) to be valid
	valid := false
	if X.board[space.i][space.j] == 0 {

		// If the current iteration is the first shift to a nearest neighbour
		// firstShift = true
		firstShift := true 
		loop:
			for {
				
				// Move space (Position) in direction of iDir and jDir
				space = Position{space.i + iDir, space.j + jDir} 
				switch {

					// If the space is out of range - Invalid
					case X.inRange(space)==false:
						break loop

					// If the space is empty - Invalid
					case X.board[space.i][space.j] == 0:
						break loop

					// If the current space is same colour - Invalid
					// This condition only applies when current space is the 
					// nearest (first) neighbour of the original starting Position
					case X.board[space.i][space.j] == X.turn && firstShift==true:
						break loop

					// If the current space is same colour, but shifted before,
					// it means that previous space is different colour, 
					// End the check 
					// This direction is considered valid
					// All pieces in this direction up to this piece should be flipped
					case X.board[space.i][space.j]==X.turn && firstShift==false:
						valid = true
						break loop

					// If the current space is a different colour, 
					// continue and shift to next space
					// The next space will not be the nearest neighbour anymore, 
					// So firstShift is set to false
					case X.board[space.i][space.j] == -X.turn:
						firstShift = false
						continue
				}
			}
		return valid
	} else { 
		return valid //false
	}
}

func (X *Board) checkValid(space Position) bool {
	// Check all directions to see if at least one direction gives a valid move
	// Returns true if for a given Position, there exists at least
	// one direction with a valid move
	valid := false 
	for _, dir := range Directions{
		if X.checkValidDir(dir.i, dir.j, space) {
			valid = true
			break
		} else {continue}
	}
	return valid
}

func (X *Board) getAllValid() []Position {
	// Obtain all valid Positions in a given Board
	validSpace := []Position{}
	for _, n := range X.neighbours {
		if X.checkValid(n) {
			validSpace = append(validSpace, Position{n.i, n.j})
		}
	}
	return validSpace
}

func (X *Board) showAllValid() {
	// Prints out whole board with valid spaces marked as 7
	Y := X.board
	for _, n := range X.neighbours {
		if X.checkValid(n) {
			Y[n.i][n.j] = 7
		}
	}
	for i := 0; i < X.length ; i++ { 
		fmt.Println(Y[i])
	}
}

func (X *Board) Move(piece Position) {
	// Place a piece on the Position (piece) given
	// Flips all relevant pieces on the board
	// Updates scores and changes turn to the next player
	if posInSlice(piece, X.validSpace) == false {
		fmt.Println("This is an invalid move")
	} else {
		flippedCount := 0
		for _, dir := range Directions{
			if X.checkValidDir(dir.i, dir.j, piece) {

				// Flip all opposing pieces in the direction until same colour is met. 
				// Given that the direction is valid
				nextPiece := Position{piece.i + dir.i, piece.j + dir.j}
				loop:
					for {

						// Move space (position) in direction of iDir and jDir
						// If next space is the same colour, flipping stops
						if X.board[nextPiece.i][nextPiece.j] == X.turn {
							break loop
						} else {
						
							// Flip the next piece and move one space in the given direction
							X.board[nextPiece.i][nextPiece.j] = X.turn
							flippedCount += 1
							nextPiece.i += dir.i
							nextPiece.j += dir.j
						}
					}
			} else {continue}
		}
		X.board[piece.i][piece.j] = X.turn // Place the piece after flipping

		// Update score
		// Total score increase = all flipped pieces + 1 new piece placed
		if X.turn == 1 {
			X.blackScore += flippedCount + 1 
			X.whiteScore -= flippedCount 
		} else {
			X.blackScore -= flippedCount
			X.whiteScore += flippedCount + 1 
		}

		// Update neighbours
		// First create empty set of neighbours and iterate
		newNeighbourSet := []Position{}
		tempNeighbourSet := append(X.neighbours, X.getNeighbour(piece)...)
		for _, n := range tempNeighbourSet {
			if n == piece {
				
				// Don't add the piece placed as a neighbour 
				continue 
			} else if posInSlice(n, newNeighbourSet) == false {
				newNeighbourSet = append(newNeighbourSet, n)
			} else {continue}
			X.neighbours = newNeighbourSet
		}	
		X.turn = -X.turn // Next player's turn
		X.validSpace = X.getAllValid() // Update valid space for next turn

		// If there are no valid moves for next player
		// Skip their turn
		if len(X.validSpace) ==0 {
			X.turn = -X.turn
			X.validSpace = X.getAllValid()

			// If there are no more moves for both players
			// End the game
			if len(X.validSpace) == 0 {
				
				// Determine Winner
				if X.blackScore > X.whiteScore {
					X.winner = 1
				} else if X.blackScore < X.whiteScore {
					X.winner = -1
				} else {
					X.winner = 99 //Draw case
				}
			}
		}
	}
}

func SetGame(state GameState) Board {
	// Setup the board for a given game state of 8x8 reversi
	// Used to restore game state from API
	// Returns a Board
	Grid := [8][8]int{}
	for i := 0; i < len(state.BlackFilled); i++ {
		Grid[state.BlackFilled[i][0]][state.BlackFilled[i][1]] = 1
	}
	for j := 0; j < len(state.WhiteFilled); j++ {
		Grid[state.WhiteFilled[j][0]][state.WhiteFilled[j][1]] = -1
	}
	B := Board{
		length:8,
		board:Grid,
		filled:[]Position{},
		empty:[]Position{},
		neighbours:[]Position{},
		validSpace:[]Position{},
		blackScore:0,
		whiteScore:0,
		winner:0,
		turn:state.Turn,
	}
	B.Setup()

	return B
}

func newGame() Board {
	// Setup the board for a new game of 8x8 reversi.
	// Returns a Board
	Grid := [8][8]int{}
	Grid[3][3] = -1
	Grid[4][4] = -1
	Grid[3][4] = 1
	Grid[4][3] = 1
	B := Board{
		length:8,
		board:Grid,
		filled:[]Position{},
		empty:[]Position{},
		neighbours:[]Position{},
		validSpace:[]Position{},
		blackScore:0,
		whiteScore:0,
		winner:0,
		turn:1,
	}
	B.Setup()

	return B
}

func simRand(game Board) Board {
	// Given a Board, simulate all moves randomly until end of game
	for {
		if game.winner == 0 {
			
			//rand is deterministic. Need to set seed 
			rand.Seed(time.Now().UTC().UnixNano()) 
			move := game.validSpace[rand.Intn(len(game.validSpace))]
			game.Move(move)		
		} else {break}
	}
	return game
}

func simRandPlus(game Board) Board {
	// Given a Board, simulate all moves semi-randomly until end of game
	// Added heuristic to discourage making very bad positions during rollouts
	// Both sides in simulation will avoid very bad positions
	// Alternative to default simRand function
	for {
		if game.winner == 0 {
			
			//rand is deterministic. Need to set seed 
			rand.Seed(time.Now().UTC().UnixNano()) 
			move := game.validSpace[rand.Intn(len(game.validSpace))]

			// When a very bad position is chosen,
			// Choose again, repeat again if very bad position chosen
			if posInSlice(move, veryBadPositions) == true {
				move = game.validSpace[rand.Intn(len(game.validSpace))]
				if posInSlice(move, veryBadPositions) == true {
					move = game.validSpace[rand.Intn(len(game.validSpace))]
				}
			}
			game.Move(move)		
		} else {break}
	}
	return game
}


func Rollout(game Board, nSim int) (int, int, int, time.Duration) {
	// Rollout function simulates nSim number of games based on given board situation
	// Function returns number of games won by black (1), white (-1), and draws and time elapsed for the function call
	turn := game.turn
	wins := 0
	draws := 0
	tempGame := game
	start := time.Now()
	for i := 0; i < nSim ; i++ {
		tempGame = simRand(game)
		if tempGame.winner == turn {
			wins++
		} 
		if tempGame.winner == 99  {
			draws++
		}
	}
	elapsed := time.Since(start)
	loss := nSim - wins - draws

	return wins, loss, draws, elapsed
}

// MCTS CODE

type Node struct {

	position Position     // Position evaluated at node
	state Board           // State of Board after position is evaluated 
	parent *Node          // Parent of node
	children []*Node      // Slice of children Nodes
	played int            // No. of times node was visited
	wins int              // No. of times won / score
	depth int             // Depth of tree - root is 0

	mobility float64      // Raw Mobility score:
						  //     The Accumulated count of valid space from explored nodes 
						  //     divided by total pieces
						  
	mobilityDenom float64 // Denominator for mobility score:
	                      //     4 x current opponent pieces 
	                      //     1 opponent piece can only have max 4 valid spaces to flip
}

func (n *Node) expandNode() {
	// Function to expand node to have children
	// Takes in validSpace array of positions from Board
	// Updates current Node
	children := []*Node{}
	for i:=0; i < len(n.state.validSpace); i++ {
		gameState := n.state // Deep copy of parent game state
		gameState.Move(n.state.validSpace[i])
		children = append(
			children, 
			&Node {
				state: gameState,
				position: n.state.validSpace[i], 
				wins: 0, 
				depth: n.depth+1,
				parent: n,
			},
		)
	}
	n.children = children
}  

func UCT(w, n, N, c int) float64 {
	// The Upper Confidence Bound  applied to Trees
	uct := float64(w)/float64(n+1) +
	       math.Sqrt(float64(c)) * math.Sqrt(math.Log(float64(N+1))/float64(n+1)) 
	return uct
}

func (n *Node) selectChild(N int, best string) *Node {
	// Selection phase for agent to choose node 
	// and decide on which Position to move
	// N = # of games played overall
	// Node selction based on upper confidence bound UCT
	index_best_score := 0
	best_uctScore := -0.00
	totalUCTScore := 0.00
	uctScore := 0.00
	// var uctScore float64

	if best == "max" {
		best_uctScore = -9999.00
	}
	if best == "min" {
		best_uctScore = 9999.0
	}
	for i, child := range n.children {
		uctScore = UCT(child.wins, child.played, N, 3)  
		
		// Adjustment score for mobility
		// Greater mobility translates to more available moves to make 
		// at later turns
		// mobScore := 0.5*math.Log(child.mobility+1)/float64(child.played +1) 

		// Adjustment score to favor inner pieces 
		// Inner pieces, or pieces close to the center of the board
		// have a higher value as they allow for more connections
		// to all other parts of the board
		innerScore := uctScore * 0.8/math.Sqrt((math.Pow((float64(child.position.i)-3.5),2)+math.Pow((float64(child.position.j)-3.5),2)))

		// Penalty for greed 
		// Squared denominator penalizes early game greed more heavily
		// Flipping more pieces early in the game is generally a bad strategy
		// Greediness gives less mobility in early to mid-games
		greedPenalty := 0.00
		if n.state.turn == 1 {
			greedPenalty = float64(child.state.blackScore - n.state.blackScore)/math.Pow(float64(child.state.blackScore + child.state.whiteScore),1)
		} 
		if n.state.turn == -1 {
			greedPenalty = float64(child.state.whiteScore - n.state.whiteScore)/math.Pow(float64(child.state.blackScore + child.state.whiteScore),1)
		}

		// Adjustment score to account for generally good / bad positions
		// - Encourages making moves that are corners
		// - Discourages making moves that give away corners
		positionScore := 0.00
		for _, corner := range corners {
			if child.position == corner {
				positionScore = uctScore *0.5
			}
		}
		for _, badpos := range badPositions {
			if child.position == badpos {
				positionScore = uctScore * (-0.15)
			}
		}
		for _, badpos := range veryBadPositions {
			if child.position == badpos {
				positionScore = uctScore * (-0.35)
			}
		}
		if best == "max" {
			totalUCTScore = uctScore  + innerScore + positionScore - greedPenalty
			totalUCTScore = uctScore 
			// totalUCTScore = uctScore + mobScore + innerScore + positionScore - greedPenalty
			if totalUCTScore  > best_uctScore {
				best_uctScore = totalUCTScore
				index_best_score = i
			}			
		} 
		if best =="min" {

			// If selecting on minimum, only select nodes that been explored 
			totalUCTScore = uctScore  - innerScore - positionScore + greedPenalty
			totalUCTScore = uctScore
			// totalUCTScore = uctScore + mobScore - innerScore - positionScore + greedPenalty
			if totalUCTScore  < best_uctScore && child.played > 0 {
				best_uctScore = totalUCTScore
				index_best_score = i
			}						
		}
	}

	return n.children[index_best_score]

	// Random method
	// rand.Seed(time.Now().UTC().UnixNano())
	// return n.children[rand.Intn(len(n.children))]
}

func backProp(n *Node, wins int, loss int, played int ) {

	// Backpropagation to traverse from child to parent nodes
	// Update count of wins and played games starting from Node n
	turn := n.state.turn // which are the wins referring to: (black:1, white:-1)
	mobility := float64(len(n.state.validSpace))
	for {
		if n.state.turn == turn {
			n.wins += wins
			n.mobility += mobility
		} else {
			n.wins += loss
		}
		n.played += played
		if n.parent == nil {
			break
		} else {
			n = n.parent
		}
	}
}

func Search(root Node, nSims int, max_iter int) Position {

	// Main function of agent to search for the optimal move
	// Expands children nodes and traverses down the tree to leaf node
	// Simulates games and backpropagates results
	// Across max_iter iterations
	// After which, selects the next move based on selectChild function
	N := 0
	wins := 0
	loss := 0
	// minScore := 999.9 // Any val greter than 1
	decision := Position{0,0}
	root.expandNode() 
	currentNode := root.selectChild(N, "min")
	wins, loss, _, _ = Rollout(currentNode.state, nSims)
	backProp(currentNode, wins, loss, nSims)
	N += nSims // Update total number of simulations 

	// Each iteration traverses down the tree
	// From parent to leaf node
	// Path of selections based on selectChild function
	// Once leaf node is reached, commence Rollout to simulate games
	for iter:=0; iter < max_iter; iter++ {
		
		// Keep selecting child nodes until leaf node is reached.		
		currentNode = root.selectChild(N, "max")
		for {
			if len(currentNode.children) == 0 { 
				break 
			} else {
				currentNode = currentNode.selectChild(N, "max")
			}
		}

		// Leaf node is reached - currentNode is leaf node
		if currentNode.played == 0 {

			// If no games played yet on this node -> rollout
			// Then backpropagate results
			wins, loss, _, _ = Rollout(currentNode.state, nSims)
			N += nSims
			backProp(currentNode, wins, loss, nSims)
			N += nSims
		} else {

			// When leaf node has been simulated before
			// Expand and look for children
			currentNode.expandNode()
			if len(currentNode.children) == 0 {

				// If there are no more children left
				// Simulate currentNode again and backpropagate				
				wins, loss, _, _ = Rollout(currentNode.state, nSims)
				N += nSims
				backProp(currentNode, wins, loss, nSims)
				N += nSims

			} else {

				// If expansion yields children, 
				// Select a child and commence rollout on child node
				// Backpropate from child node
				currentNode = currentNode.selectChild(N, "max")
				wins, loss, _, _ = Rollout(currentNode.state, nSims)
				N += nSims
				backProp(currentNode, wins, loss, nSims)
				N += nSims
			}
		}
	}

	// Once all simulation and max iterations reached
	// Select the child from the root node
	// This will be the move the agent makes
	decision = root.selectChild(0, "min").position

	return decision
}


// func main() {
// 	// time_start := time.Now()
// 	game := newGame()
// 	move := Position{0,0}
// 	root := Node{
// 		// state:newGame(),
// 		state:game,
// 		depth:0,
// 	}
// 	move = Search(root, 1, 5000)
// 	fmt.Println("$$$ ", move, game.turn)



// 	// move := Position{0,0}
// 	// root := Node{
// 	// 	// state:newGame(),
// 	// 	state:game,
// 	// 	depth:0,
// 	// }
// 	// move = Search(root, 1, 5000)
// 	// 	fmt.Println("$$$ ", move, game.turn)
// 	// game.Show()
// 	// fmt.Println(game.blackScore, game.whiteScore)
// 	// elapsed := time.Since(time_start)
// 	// fmt.Println("Time elapsed: ", elapsed)


// 	//####

func Simulator(N int, nSims int, max_iter int) {
	// Simulate N games with agent pitted against random play
	// Returns # of games won/lost/draw for MCTS agent
	// Used as benchmark testing against Search() function 
	nBlackWins :=0
	nWhiteWins :=0
	nDraws :=0

	for nGames:=0; nGames < N; nGames++ {
		game := newGame()
		root := Node{
			state:game,
			depth:0,
		}
		move := Position{0,0}
		// Black plays as MCTS, white random
		for {
			if game.winner == 0 {
				if game.turn == 1 {
					root = Node{state:game, depth:0}
					move = Search(root, nSims, max_iter)
					game.Move(move)					
					// fmt.Println("Black moves: ", move.PrintPrettifyNotation(), game.blackScore, game.whiteScore)
				} else {
					rand.Seed(time.Now().UTC().UnixNano()) 
					move = game.validSpace[rand.Intn(len(game.validSpace))]
					// When a very bad position is chosen,
					// Choose again, repeat again if very bad position chosen
					if posInSlice(move, veryBadPositions) == true {
						move = game.validSpace[rand.Intn(len(game.validSpace))]
						if posInSlice(move, veryBadPositions) == true {
							move = game.validSpace[rand.Intn(len(game.validSpace))]
						}
					}
					// fmt.Println("White moves: ", move.PrintPrettifyNotation(), game.blackScore, game.whiteScore)
					game.Move(move)	
				}
			} else {break}
		}
		fmt.Println("Game #", nGames, game.winner, game.blackScore, game.whiteScore)
		if game.winner == 1 {
			nBlackWins +=1
		}
		if game.winner == -1 {
			nWhiteWins +=1
		}
		if game.winner == 99 {
			nDraws +=1
		}
	}
	fmt.Printf("Black Wins: %d , White Wins: %d , Draws: %d \n", nBlackWins, nWhiteWins, nDraws)
}


func RandomRandomPlay(N int) {
	// Simulate N games with random play pitted against random play
	// Returns # of games won/lost/draw for Black vs White
	nBlackWins :=0
	nWhiteWins :=0
	nDraws :=0
	for nGames:=0; nGames < N; nGames++ {
		game := newGame()
		move := Position{0,0}
		for {
			if game.winner == 0 {
				if game.turn == 1 {
					rand.Seed(time.Now().UTC().UnixNano()) 
					move = game.validSpace[rand.Intn(len(game.validSpace))]
					// When a very bad position is chosen,
					// Choose again, repeat again if very bad position chosen
					if posInSlice(move, veryBadPositions) == true {
						move = game.validSpace[rand.Intn(len(game.validSpace))]
						if posInSlice(move, veryBadPositions) == true {
							move = game.validSpace[rand.Intn(len(game.validSpace))]
						}
					}
				} else {
					rand.Seed(time.Now().UTC().UnixNano()) 
					move = game.validSpace[rand.Intn(len(game.validSpace))]
					// When a very bad position is chosen,
					// Choose again, repeat again if very bad position chosen
					if posInSlice(move, veryBadPositions) == true {
						move = game.validSpace[rand.Intn(len(game.validSpace))]
						if posInSlice(move, veryBadPositions) == true {
							move = game.validSpace[rand.Intn(len(game.validSpace))]
						}
					}
				}
				// fmt.Println("White moves: ", move.PrintPrettifyNotation(), game.blackScore, game.whiteScore)
				game.Move(move)	
			} else {break}
		}
		fmt.Println("Game #", nGames, game.winner, game.blackScore, game.whiteScore)
		if game.winner == 1 {
			nBlackWins +=1
		}
		if game.winner == -1 {
			nWhiteWins +=1
		}
		if game.winner == 99 {
			nDraws +=1
		}
	}
	fmt.Printf("Black Wins: %d , White Wins: %d , Draws: %d \n", nBlackWins, nWhiteWins, nDraws)
}