package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Reversi, %q", html.EscapeString(r.URL.Path))
	http.ServeFile(w, r, "./static/index.html")
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
		state: game,
		depth: 0,
	}
	decision := Search(root, 20, 300)
	game.Move(decision)

	response := DecisionResponse{
		Move:       [2]int{decision.i, decision.j},
		Colour:     state.Turn,
		Turn:       game.turn,
		BlackScore: game.blackScore,
		WhiteScore: game.whiteScore,
	}

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)

}
