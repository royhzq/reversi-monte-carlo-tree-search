package main

import (
	"fmt"
	"time"
	"math"
	"strconv"
	"math/rand"
)

type Tuple struct {
	i int
	j int
}

type strTuple struct {
	i string
	j string
}

func (tuple Tuple) PrintPrettifyNotation() strTuple {
	alphabet := []string{"A","B","C","D","E","F","G","H"}
	return strTuple{alphabet[tuple.j], strconv.Itoa(tuple.i+1)}
}

// func (strtuple strTuple) convertTuple() Tuple {
// 	alphabet := []string{"A","B","C","D","E","F","G","H"}
// 	for j, letter := range alphabet {
// 		if letter == strtuple.i
// 	}	
// 	strtuple.
// }

type Board struct {
	
	length int //max length of board (i.e., 8 for standard board size)
	board [8][8]int //state of board
	filled []Tuple //Slice of all spaces filled up by a piece
	empty []Tuple //Slice of all spaces that remain empty
	neighbours []Tuple //Slice of all empty neighbours of pieces
	validSpace []Tuple //Slice of all valid moves for the player turn
	blackScore int //Total number of Black (1)  
	whiteScore int //Total number of White (-1) 
	winner int //Winner of game - Black (1), White (-1), Draw (99), Undetermined (0). Undetermined is default
	turn int //Whose turn is it (1 for Black, -1 for White)
}

func tupleInSlice(a Tuple, list []Tuple) bool {
	//General function that checks if a Tuple is in a slice
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func (X Board) Show() {
	//Method to show board on command line
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

func (X Board) inRange(space Tuple) bool {
	// Check if tuple is inside the board (i.e., not out of range)
	// Return false if Tuple is out of range
	if space.i<X.length && space.j<X.length && space.i>=0 && space.j>=0 {
		//Index cannot be more than length or board. Also cannot be negative
		return true
	} else {
		return false
	}

} 

func (X *Board) getNeighbour(piece Tuple) []Tuple {
	//Get immediate neighbouring empty spaces of a given coordinate (tuple)
	//There are 
	
	neighbours := []Tuple{}

	if X.inRange(Tuple{piece.i+1, piece.j}) && X.board[piece.i+1][piece.j]==0 {
		neighbours = append(neighbours, Tuple{piece.i+1,piece.j})
	}
	if X.inRange(Tuple{piece.i+1, piece.j+1}) && X.board[piece.i+1][piece.j+1]==0 {
		neighbours = append(neighbours, Tuple{piece.i+1,piece.j+1})
	}
	if X.inRange(Tuple{piece.i, piece.j+1}) && X.board[piece.i][piece.j+1]==0 {
		neighbours = append(neighbours, Tuple{piece.i,piece.j+1})
	}
	if X.inRange(Tuple{piece.i-1, piece.j+1}) && X.board[piece.i-1][piece.j+1]==0 {
		neighbours = append(neighbours, Tuple{piece.i-1,piece.j+1})
	}
	if X.inRange(Tuple{piece.i-1, piece.j})  && X.board[piece.i-1][piece.j]==0 {
		neighbours = append(neighbours, Tuple{piece.i-1,piece.j})
	}
	if X.inRange(Tuple{piece.i-1, piece.j-1}) && X.board[piece.i-1][piece.j-1]==0 {
		neighbours = append(neighbours, Tuple{piece.i-1,piece.j-1})
	}
	if X.inRange(Tuple{piece.i, piece.j-1}) && X.board[piece.i][piece.j-1]==0 {
		neighbours = append(neighbours, Tuple{piece.i,piece.j-1})
	}
	if X.inRange(Tuple{piece.i+1, piece.j-1}) && X.board[piece.i+1][piece.j-1]==0 {
		neighbours = append(neighbours, Tuple{piece.i+1,piece.j-1})
	}
	
	return neighbours
}

func (X *Board) initNeighbours() {
	//Method that initializes all neighbours on board
	filled := X.filled 
	neighbours := []Tuple{} //List to store all neighbours
	pieceNeighbour := []Tuple{} //List to store neighbours of a single piece
	Y := X.board //Temp board

	//Loop through all filled spaces and get their neighbours
	//Add neighbours to neighbours list if tuple does not already exist in list
	for _, piece := range filled {
		pieceNeighbour = X.getNeighbour(piece)
		for _, n := range pieceNeighbour {

			//if neighbour tuple is not in neighbours slice. Add to slice
			if tupleInSlice(n, neighbours)==false {
				neighbours = append(neighbours, n)
				Y[n.i][n.j] = 5
			}
		}

	}

	X.neighbours = neighbours
}

func (X *Board) getScores() (int, int) {
	// Calculates overall score for (black, white)
	// Black is 1 and White is -1
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
	//Given a Board, initialize parameters 
	//Run this method one time after initializing Board

	X.blackScore = 0
	X.whiteScore = 0
	X.filled = []Tuple{}
	X.empty = []Tuple{}
	X.neighbours = []Tuple{}
	X.validSpace = []Tuple{}

	//Initialize empty set
	for i := 0; i < X.length; i++ { 
		for j := 0; j < X.length; j++ {
			if X.board[i][j] ==0 {
				X.empty = append(X.empty, Tuple{i,j})
			}
		}
	}

	//Initialize filled set
	for i := 0; i < X.length; i++ { 
		for j := 0; j < X.length; j++ {
			if X.board[i][j] !=0 {
				X.filled = append(X.filled, Tuple{i,j})
			}
		}
	}


	//Initialize set of all neighbours at the start
	X.initNeighbours()

	//Initialize validSpace - Must be done after neighbours initialized
	X.validSpace = X.getAllValid()

	//Initialize beginning scores 
	X.blackScore, X.whiteScore = X.getScores()

	//Initialize winner - Undetermined (0)
	X.winner = 0

}

func (X *Board) checkValidDir(iDir int, jDir int, space Tuple) bool {
	
	// Check, in the direction (iDir, jDir), whether pieces will flip
	// Direction is defined by (iDir, jDir), for example,
	// if (iDir=1, jDir=1), we are checking the N-E direction of a given space (Tuple)
	// if (iDir=0, jDir=-1), we are checking the West direction of a given space (Tuple)

	//First Check if space is empty. Space has to be empty (0) to be valid

	valid := false

	if X.board[space.i][space.j]==0 {
		firstShift := true //Shift to first neighbour

		loop:
			for {

				space = Tuple{space.i+iDir, space.j+jDir} //Move space (tuple) in direction of iDir and jDir
				switch {

					//If the space is out of range - Invalid
					case X.inRange(space)==false:
						// fmt.Println("case 1", space)
						break loop

					//If the space is empty - Invalid
					case X.board[space.i][space.j] == 0:
						// fmt.Println("case 2", space)
						break loop

					//If the current space is same colour - Invalid
					//This only applies when current space is the first neighbour of the original starting space
					case X.board[space.i][space.j] == X.turn && firstShift==true:
						// fmt.Println("case 3", space)
						break loop

					//If the current space is same colour, but shifted before,
					//means that previous space is different colour, end the check and this
					//is where all pieces in this direction up to this piece is flipped
					//this will then be valid
					case X.board[space.i][space.j]==X.turn && firstShift==false:
						valid = true
						// fmt.Println(space.i,space.j,"TRUE")
						break loop

					//If the current space is differing colour - Continue and shift to next space
					//Will not be the first shift anymore, so firstShift = false
					case X.board[space.i][space.j] == -X.turn:
						firstShift = false
						// fmt.Println(space.i,space.j,"Moving On..")
						continue
				}
				
			}
		return valid

	} else { 
		return valid //false
	}
	
}

func (X *Board) checkValid(space Tuple) bool {
	//Given a space, check all directions to see if at least one direction gives a valid move
	//directions start from North, going clockwise
	directions := []Tuple{
		Tuple{-1,0},
		Tuple{-1,1},
		Tuple{0,1},
		Tuple{1,1},
		Tuple{1,0},
		Tuple{1,-1},
		Tuple{0,-1},
		Tuple{-1,-1},
	}

	valid := false 

	for _, dir := range directions{
		if X.checkValidDir(dir.i, dir.j, space) {
			valid = true
			break
		} else {continue}

	}
	return valid
}

func (X *Board) getAllValid() []Tuple {

	validSpace := []Tuple{}
	for _, n := range X.neighbours {
		if X.checkValid(n) {
			validSpace = append(validSpace, Tuple{n.i,n.j})
		}
	}
	return validSpace
}


func (X *Board) showAllValid() {
	//prints out whole board with valid spaces marked as 7
	//loops through all neighbours to check validity
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

func (X *Board) Move(piece Tuple) {

	if tupleInSlice(piece, X.validSpace)==false {
		fmt.Println("This is an invalid move")
	} else {

		directions := []Tuple{
			Tuple{-1,0},
			Tuple{-1,1},
			Tuple{0,1},
			Tuple{1,1},
			Tuple{1,0},
			Tuple{1,-1},
			Tuple{0,-1},
			Tuple{-1,-1},	
		}

		flippedCount := 0

		for _, dir := range directions{

			if X.checkValidDir(dir.i, dir.j, piece) {
				// Validity has been checked 
				// Flip all opposing pieces in the direction until same colour is met. 

				nextPiece := Tuple{piece.i+dir.i, piece.j+dir.j}
				loop:
					for {
						 //Move space (tuple) in direction of iDir and jDir
						if X.board[nextPiece.i][nextPiece.j] == X.turn {
							break loop
						} else {
							//Flip the next piece and move one space in the direction
							X.board[nextPiece.i][nextPiece.j] = X.turn
							flippedCount += 1
							nextPiece.i += dir.i
							nextPiece.j += dir.j
						}
					}
			} else {continue}

		}

		X.board[piece.i][piece.j] = X.turn // Place the piece only after flipping

		//Update score
		if X.turn==1 {
			X.blackScore += flippedCount + 1 // Total score increase is = all flipped pieces + 1 new piece placed
			X.whiteScore -= flippedCount 
		} else {
			X.blackScore -= flippedCount
			X.whiteScore += flippedCount + 1 // Total score increase is = all flipped pieces + 1 new piece placed
		}

		//Update neighbours
		//First create empty set of neighbours and loop through..

		newNeighbourSet := []Tuple{}
		tempNeighbourSet := append(X.neighbours, X.getNeighbour(piece)...)
		for _, n := range tempNeighbourSet {

			if n == piece {
				continue // Don't add the piece as a neighbour 
			} else if tupleInSlice(n, newNeighbourSet) == false {

				// If not in neighbourset, add to new neighbourset
				newNeighbourSet = append(newNeighbourSet, n)
			} else {continue}

		X.neighbours = newNeighbourSet

		}

		
		X.turn= -X.turn //Next player turn
		X.validSpace = X.getAllValid() //Update valid space for next turn
		//fmt.Println(X.validSpace, len(X.validSpace))

		// - If no valid moves, skip
		if len(X.validSpace) ==0 {
			//fmt.Println(X.turn, " Has No valid moves. Skip a turn.")
			X.turn= -X.turn
			X.validSpace = X.getAllValid()

			// No more moves for both players. Game ended.
			if len(X.validSpace) == 0 {
				//fmt.Println(X.turn, "Also Has No valid moves. Game has ended.")
				
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
		filled:[]Tuple{},
		empty:[]Tuple{},
		neighbours:[]Tuple{},
		validSpace:[]Tuple{},
		blackScore:0,
		whiteScore:0,
		winner:0,
		turn:state.Turn,
	}
	B.Setup()
	return B
}

func newGame() Board {

	//Setup the board for a new game of 8x8 reversi.
	//Returns a Board

	Grid := [8][8]int{}
	Grid[3][3] = -1
	Grid[4][4] = -1
	Grid[3][4] = 1
	Grid[4][3] = 1

	B := Board{
		length:8,
		board:Grid,
		filled:[]Tuple{},
		empty:[]Tuple{},
		neighbours:[]Tuple{},
		validSpace:[]Tuple{},
		blackScore:0,
		whiteScore:0,
		winner:0,
		turn:1,
	}
	B.Setup()
	// B.Show()

	return B
}

func simRand(game Board) Board {
	
	//Given a board, simulate all moves randomly until end of game
	for {
		if game.winner == 0 {
			rand.Seed(time.Now().UTC().UnixNano()) //rand is deterministic. Need to set seed 
			move := game.validSpace[rand.Intn(len(game.validSpace))]
			game.Move(move)		

		} else {break}
	}
	return game

}

func Rollout(game Board, nSim int) (int, int, int, time.Duration) {

	//Rollout function simulates nSim number of games based on given board situation
	//Function returns number of games won by black (1), white (-1), and draws and time elapsed for the function call
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

	position Tuple // position evaluated at node
	state Board // state of game (Board) after position is evaluated 
	parent *Node // parent of node
	children []*Node //array of children Nodes
	played int // no. of times visited
	wins int // no. of times won / score
	depth int // depth of tree - root is 0
	mobility float64 // Raw Mobility score - Accumulated count of valid space from explored nodes divided by total pieces
	mobilityDenom float64 //Denominator for mobility score - 4 x current opponent pieces (because 1 opponent piece can only have max 4 valid spaces to eat)
}

func (n *Node) expandNode() {
	// Function to expand node to have children
	// Takes in validSpace array of tuples from Board
	// Updates current Node
	children := []*Node{}

	for i:=0; i < len(n.state.validSpace); i++ {

		gameState := n.state // Deep Copy parent game state
		gameState.Move(n.state.validSpace[i])

		children = append(
			children, 
			&Node {
				state:gameState,
				position:n.state.validSpace[i], 
				wins:0, 
				depth:n.depth+1,
				parent:n,
			},
		)
	}
	n.children = children
}  

func UCT(w, n, N, c int) float64 {

	// Calculates Upper Confidence Bound 1 applied to Trees
	uct := float64(w)/(float64(n+1)) + math.Sqrt(float64(c))*math.Sqrt(math.Log(float64(N+1))/float64(n+1)) 

	return uct
}

func (n *Node) selectChild(N int, best string) *Node {
	
	// N = # of games played overall
	// max = true means select child with highest score
	// max = false selects child with lowest score

	// Selection phase

	// UCT method
	index_best_score := 0
	
	best_uctScore:= -0.00
	totalUCTScore:= 0.00

	if best == "max" {
		best_uctScore = -9999.00
	}
	if best == "min" {
		best_uctScore = 9999.0
	}	
	var uctScore float64

	corners := []Tuple{{0,0}, {0,7}, {7,0}, {7,7}}
	badPositions := []Tuple{
		{0,1},{1,0},
		{6,0},{7,1},
		{7,6},{6,7},
		{0,6},{1,7},
	}
	veryBadPositions := []Tuple{
		{1,1},{1,6},{6,1},{6,6},
	}
	for i, child := range n.children {

		uctScore = UCT(child.wins, child.played, N, 3)  
		
		// Adjustment score for mobility
		mobScore := 0.5*math.Log(child.mobility+1)/float64(child.played +1) 

		//Adjustment score to favor inner pieces 
		innerScore := uctScore * 0.8/math.Sqrt((math.Pow((float64(child.position.i)-3.5),2)+math.Pow((float64(child.position.j)-3.5),2)))
		
		//Adjustment score to account for universally good / bad positions
		positionScore := 0.00

		//Penalty for greed. Squared denominator penalizes early game greed more heavily
		greedPenalty := 0.00
		if n.state.turn == 1 {
			greedPenalty = float64(child.state.blackScore - n.state.blackScore)/math.Pow(float64(child.state.blackScore + child.state.whiteScore),1)
			// fmt.Println(child.position ,"Change in BlackScore: ", child.state.blackScore - n.state.blackScore, child.state.blackScore + child.state.whiteScore, greedPenalty)
		} 
 	
		if n.state.turn == -1 {
			greedPenalty = float64(child.state.whiteScore - n.state.whiteScore)/math.Pow(float64(child.state.blackScore + child.state.whiteScore),1)
			// fmt.Println(child.position ,"Change in BlackScore: ", child.state.whiteScore - n.state.whiteScore, child.state.blackScore + child.state.whiteScore, greedPenalty)
		}


		for _, corner := range corners {
			if child.position == corner {
				positionScore = uctScore *0.35
			}
		}
		for _, badpos := range badPositions {
			if child.position == badpos {
				positionScore = uctScore * (-0.15)
			}
		}
		for _, badpos := range veryBadPositions {
			if child.position == badpos {
				positionScore = uctScore * (-0.65)
			}
		}
		
		if best == "max" {
			
			totalUCTScore = uctScore + mobScore + innerScore + positionScore - greedPenalty

			if totalUCTScore  > best_uctScore {
				best_uctScore = totalUCTScore
				index_best_score = i
			}			
		} 
		if best =="min" {
			// If selecting on minimum, only select nodes that been explored 
			totalUCTScore = uctScore + mobScore - innerScore - positionScore + greedPenalty
			if totalUCTScore  < best_uctScore && child.played > 0 {
				best_uctScore = totalUCTScore
				index_best_score = i
			}						
		}
		
		// fmt.Printf("{%d,%d}, UCT: %f (%.0f%%), mob: %f (%.0f%%), inner: %f (%.0f%%), position: %f (%.0f%%), greed: %f (%.0f%%), total UCT: %f (100%%), visited: %d \n",
		// 	child.position.i,child.position.j,
		// 	uctScore, uctScore/totalUCTScore *100,
		// 	mobScore, mobScore/totalUCTScore *100,			
		// 	innerScore, innerScore/totalUCTScore *100,			
		// 	positionScore, positionScore/totalUCTScore *100,			
		// 	greedPenalty, greedPenalty/totalUCTScore *100,			
		// 	totalUCTScore, 		
		// 	child.played, 
		// )
	}
	// fmt.Printf("SELECTED: {%d,%d}, Depth: %d \n", 
	// 	n.children[index_best_score].position.i ,n.children[index_best_score].position.j,  
	// 	n.children[index_best_score].depth,
	// )
	return n.children[index_best_score]

	// Random method
	// rand.Seed(time.Now().UTC().UnixNano())
	// return n.children[rand.Intn(len(n.children))]
}

func backProp(n *Node, wins int, loss int, played int ) {
	// Backpropagation to add count of wins and played games starting from Node n
	// wins : # of wins in sim , played : # of games played in sim
	// wins refer to the # of wins of the turn of current node n
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


func Search(root Node, nSims int, max_iter int) Tuple {

	N := 0
	wins := 0
	loss := 0
	// minScore := 999.9 // Any val greter than 1
	decision := Tuple{0,0}

	root.expandNode()
	currentNode := root.selectChild(N, "min")
	wins, loss, _, _ = Rollout(currentNode.state, nSims)
	backProp(currentNode, wins, loss, nSims)
	N += nSims // Update total number of simulations 

	for iter:=0; iter < max_iter; iter++ {
		
		// Keep selecting until leaf node is reached.
		// This is 1 iteration
		
		currentNode = root.selectChild(N, "max")
		for {
			if len(currentNode.children) == 0 { 
				break 
			} else {
				currentNode = currentNode.selectChild(N, "max")
				// fmt.Println("Traversing down.. Depth: ", currentNode.depth)
			}
		}

		// Leaf node is current node

		if currentNode.played == 0 {
			//If no games played on this node, rollout simulate
			wins, loss, _, _ = Rollout(currentNode.state, nSims)
			N += nSims
			backProp(currentNode, wins, loss, nSims)
			N += nSims
		} else {

			currentNode.expandNode()
			if len(currentNode.children) == 0 {
				// if there are no more children, rollout current node
				wins, loss, _, _ = Rollout(currentNode.state, nSims)
				N += nSims
				backProp(currentNode, wins, loss, nSims)
				N += nSims

			} else {
				// if there are children, select and rollout new node
				currentNode = currentNode.selectChild(N, "max")
				wins, loss, _, _ = Rollout(currentNode.state, nSims)
				N += nSims
				backProp(currentNode, wins, loss, nSims)
				N += nSims
			}

		}
	

	}

	// fmt.Println("DECISION BASED ON ORIGINAL ")
	// for _, child := range root.children {
	// 	// fmt.Println("orig: ", child.position, child.wins, child.played, float64(child.wins)/float64(child.played+1) + child.mobility)
	// 	if float64(child.wins)/float64(child.played+1) + child.mobility < minScore &&  float64(child.wins)/float64(child.played+1) + child.mobility > 0 {
	// 		minScore = float64(child.wins)/float64(child.played+1) + child.mobility
	// 		decision = child.position
	// 	}
	// }
	// fmt.Println("selectChild: ", root.selectChild(0, "min").position)
	// fmt.Println("original: ", decision, minScore)
	decision = root.selectChild(0, "min").position
	return decision
}


// func main() {
// 	// time_start := time.Now()
// 	game := newGame()
// 	move := Tuple{0,0}
// 	root := Node{
// 		// state:newGame(),
// 		state:game,
// 		depth:0,
// 	}
// 	move = Search(root, 1, 5000)
// 	fmt.Println("$$$ ", move, game.turn)

// game.Move(Tuple{2,3})
// game.Move(Tuple{2,2})
// game.Move(Tuple{3,2})
// game.Move(Tuple{2,4})
// game.Move(Tuple{3,5})
// game.Move(Tuple{4,6})
// game.Move(Tuple{5,4})
// game.Move(Tuple{4,2})
// game.Move(Tuple{3,6})
// game.Move(Tuple{5,5})
// game.Move(Tuple{5,3})
// game.Move(Tuple{3,7})
// game.Move(Tuple{2,5})
// game.Move(Tuple{5,2})
// game.Move(Tuple{4,5})
// game.Move(Tuple{5,6})
// game.Move(Tuple{1,3})
// game.Move(Tuple{0,2})
// game.Move(Tuple{3,1})
// game.Move(Tuple{4,0})
// game.Move(Tuple{6,5})
// game.Move(Tuple{7,4})
// game.Move(Tuple{1,4})
// game.Move(Tuple{1,2})
// game.Move(Tuple{4,1})
// game.Move(Tuple{5,0})
// game.Move(Tuple{6,3})
// game.Move(Tuple{7,2})
// game.Move(Tuple{7,5})
// game.Move(Tuple{0,5})
// game.Move(Tuple{6,2})
// game.Move(Tuple{5,1})
// game.Move(Tuple{6,4})
// game.Move(Tuple{7,3})
// game.Move(Tuple{2,1})
// game.Move(Tuple{1,5})
// game.Move(Tuple{2,6})
// game.Move(Tuple{3,0})
// game.Move(Tuple{5,7})
// game.Move(Tuple{2,7})
// game.Move(Tuple{2,0})
// game.Move(Tuple{4,7})
// game.Move(Tuple{0,4})
// game.Move(Tuple{0,3})
// game.Move(Tuple{7,1})
// game.Move(Tuple{1,0})
// game.Move(Tuple{1,7})
// game.Move(Tuple{6,1})
// game.Move(Tuple{7,0})
// game.Move(Tuple{1,6})







// 	// move := Tuple{0,0}
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

	// time_start := time.Now()

	// nBlackWins :=0
	// nWhiteWins :=0
	// nDraws :=0

	// for nGames:=0; nGames < 100; nGames++ {
	// 	game := newGame()
	// 	root := Node{
	// 		// state:newGame(),
	// 		state:game,
	// 		depth:0,
	// 	}
	// 	move := Tuple{0,0}
	// 	// Black plays as MCTS, white random
	// 	for {
	// 		if game.winner == 0 {
	// 			if game.turn == 1 {
	// 				root = Node{state:game, depth:0}
	// 				move = Search(root, 1, 100)
	// 				// game.Show()
	// 				game.Move(move)
	// 				// fmt.Println("*************")
	// 				// game.Show()
					
	// 				fmt.Println("Black moves: ", move.PrintPrettifyNotation(), game.blackScore, game.whiteScore)

	// 			} else {
	// 				rand.Seed(time.Now().UTC().UnixNano()) //rand is deterministic. Need to set seed 
	// 				move = game.validSpace[rand.Intn(len(game.validSpace))]
	// 				game.Move(move)	
	// 				fmt.Println("White moves: ", move.PrintPrettifyNotation(), game.blackScore, game.whiteScore)

	// 			}

	// 		} else {break}
	// 	}

	// 	// game.Show()
	// 	fmt.Println("Game # ", nGames, "**************")
	// 	fmt.Println(game.winner, game.blackScore, game.whiteScore)

	// 	if game.winner == 1 {
	// 		nBlackWins +=1
	// 	}
	// 	if game.winner == -1 {
	// 		nWhiteWins +=1
	// 	}
	// 	if game.winner == 99 {
	// 		nDraws +=1
	// 	}
	// }
	// fmt.Printf("Black Wins: %d , White Wins: %d , Draws: %d \n", nBlackWins, nWhiteWins, nDraws)




	// elapsed := time.Since(time_start)
	// fmt.Println("Time elapsed: ", elapsed)
// }
