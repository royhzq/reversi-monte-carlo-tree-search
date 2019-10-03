package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"io/ioutil"
)

type GameState struct {

	// Struct to hold game state information 
	// from JSON body posted to endpoint -> converted to GameState 

	BlackFilled [][2]int `json:"blackFilled"` // Currently filled black pieces on board 
	WhiteFilled [][2]int `json:"whiteFilled"` // Currently filled white pieces on board
	Turn        int      `json:"turn"`        // Agent's turn to play as (1 for black, -1 for white)
}

type DecisionResponse struct {

	// Struct to hold response information on move made by AI
	// in response to the gamestate posted by user to endpoint

	Move       [2]int `json:"move"`        // Decision of agent for Position to place piece
	Colour     int    `json:"colour`       // Colour of the piece placed by agent
	Turn       int    `json:"turn"`        // Whose turn it is after move is made
	BlackScore int    `json:"blackScore"`  // Black's Score after Move is made
	WhiteScore int    `json:"whiteScore"`  // White's Score after Move is made
}


func Index(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Reversi, %q", html.EscapeString(r.URL.Path))
	http.ServeFile(w, r, "./static/index.html" )
}


func GameStateAPI(w http.ResponseWriter, r *http.Request) {
	/*  API Endpoint to receive JSON gamestate from POST request 
	Request JSON example:
		{
			"blackFilled":[[3,3],[4,4]],    // Positions on board filled with black pieces
			"whiteFilled":[[3,4],[4,3]],    // Positions on board filled with white piece
			"turn":1,                       // Agent's turn to play as (1 black, -1 white)
		}
	Response JSON example:
		{
			"move":[3,2],                   // The move the agent is going to make
			"turn":1,                       // The turn
			"blackScore":1,
			"whiteScore":4
		}
	*/
	var state = GameState{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &state) //Converts request body to moveArr type
	if err != nil {
		fmt.Println("error:", err)
	}
	game := SetGame(state)
	root := Node{
		state:game,
		depth:0,
	}
	decision := Search(root, 20, 300) 
	game.Move(decision)
	
	response := DecisionResponse{
		Move:[2]int{decision.i, decision.j},
		Colour:state.Turn,
		Turn:game.turn,
		BlackScore:game.blackScore,
		WhiteScore:game.whiteScore,
	}

	//Allow CORS here By * or specific origin
    w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)

}