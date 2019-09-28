
# About reversi-mcts

reversi-mcts is a reversi/otherllo agent that is based on the [Monte Carlo Tree Search (MCTS)](https://en.wikipedia.org/wiki/Monte_Carlo_tree_search) algorithm implemented in the [Go](https://golang.org) programming language.

It includes both an API client and a javascript based interface to play against the agent.

Give it a go and test your skills against this agent at [royhung.com/reversi](https://royhung.com/reversi)

# Installing reversi-mcts

Make sure you have a working Go installation, see the installation guide at http://golang.org/doc/install.html 

Next, make sure ``` $GOPATH ``` is set (e.g. as ~/go) and:

```console
$ cd $GOPATH/src
$ git clone https://gitlab.com/royhung_/reversi-monte-carlo-tree-search.git
$ cd reversi-monte-carlo-tree-search
$ go build
```

# Using reversi-mcts

Run the built package to start the server:

```console
$ ./reversi-monte-carlo-tree-search

```

The application will be running on http://localhost:8080 with a reversi board interface and a playable MCTS agent. 


# API Endpoint

To access the agent's API directly, make a POST request to ```/search_move```. <br>The endpoint consumes data on the state of the game via the positions of black pieces and white pieces in play. <br>
The response from the MCTS agent will contain the coordinates of the position the agent will place its next move. 

``` POST /search_move ```

### Request POST JSON Example

```json
{
    "blackFilled":[[3,3],[4,4]],    
    "whiteFilled":[[3,4],[4,3]],    
    "turn":1,                      
}
```
| Property | Type |Description |
| --- | --- | :- |
| ``` blackFilled ``` | Object | Array of coordinate positions [ i , j ] of black pieces, where i refers to the ith row on board and j refers to the jth row on the board   |
| ``` whiteFilled ``` | Object | Array of coordinate positions [ i , j ] of white pieces, where i refers to the ith row on board and j refers to the jth row on the board  |
| ``` turn ``` | Integer | The colour agent is supposed to play as for its turn (1 black, -1 white) |


### Response POST JSON Example

```json
{
    "move":[3,2],                   
    "turn":1,                       
    "blackScore":1,
    "whiteScore":4
}
```
| Property | Type |Description |
| --- | --- | :- |
| ``` move ``` | Object | Array containing coordinate position [ i, j ] of the move the agent has made   |
| ``` turn ``` | Integer | The colour of the piece placed (1 black, -1 white) |
| ``` blackScore ``` | Integer | The resulting number of black pieces on the board after move is made |
| ``` whiteScore ``` | Integer | The resulting number of white pieces on the board after move is made |


# More information

For a more detailed write up on the algorithm, performance and the parameters used. Please visit https://royhung.com/reversi
