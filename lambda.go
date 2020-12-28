// Alternative compilation for deploying as an AWS lambda function
// > go build lambda.go reversi.go
// > zip function.zip lambda
// Upload to AWS Lambda console accordingly

package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleLambdaEvent(state GameState) (DecisionResponse, error) {

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

	return response, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
