/*!
 * Reversi Game Board
 * 2019 Roy Hung
 */

function updateMessage(text) {
    /**
    * Changes text in 'message' element
    * Updates status of game for user
    * @param {string} text - The text to display
    */
    var msg = document.getElementById("message");
    msg.innerHTML = text;
}

function turn2Colour(turn) {
    /**
    * Converts integer Converts integer 1, -1 into representative colour
    * 1 -> 'Black'
    * -1 -> 'White'
    * @param {number} turn - Integer 1 or -1
    * @returns {string} The current player's turn colour (black or white) 
    */
    if (turn === 1) {
        return 'Black';
    } else if (turn === -1) {
        return 'White';
    } else {
        return '';
    }
}

function inArray(target, array) {
    /** 
    * To check if target element is contained in array
    * @param {Object} target - Element to check if it is array
    * @param {Array} array - Array of elements to check against
    * @returns {boolean} true if target element is in array, else falses
    */
    for(var i = 0; i < array.length; i++) {
        if(array[i] === target) {
            return true
        }
    }
    return false;
}

function newBoard(dim) {
    /** 
    * Creates a 2D array representing the reversi board
    * Initializes starting pieces for a new game
    * Returns dim x dim 2D array
    * @param {number} dim - dimension of square board (i.e. 8x8)
    * @returns {Array} board - An array of arrays representing the 2D reversi board
    */
    var board =[];
    for (var i=0; i <dim; i++) {
        board.push(Array(dim))
    }
    // Setup starting pieces
    var k = dim/2 -1;
    board[k][k] = -1;
    board[k][k+1] = 1;
    board[k+1][k] = 1;
    board[k+1][k+1] = -1;

    return board;
}

function drawGrid(width, dim) {
    /** 
    * Render the reversi board in an svg element
    * @param {number} width - pixel width for board size
    * @param {number} dim - dimension of square board (i.e. 8x8)
    */
    var svg = document.getElementById("board");
    svg.setAttribute('data-turn', '1');
    svg.setAttribute('data-filled', '[]');
    svg.setAttribute('data-neighbours', '[]');
    svg.setAttribute('data-validspace', '[]');
    svg.setAttribute('data-blackscore', '');
    svg.setAttribute('data-whitescore', '');
    sqLength = width/dim; // Initialize size of each square
    // Draw border for the reversi board
    gridBorder = ('<rect class="gridBorder" '
        + 'x=0 ' 
        + 'y=0 '
        + 'width=' + width + ' '
        + 'height=' + width + ' '
        + 'style="fill:none; stroke-width:1 ; stroke:rgb(0,0,0)" />'
    );
    svg.innerHTML += gridBorder;
    // Loop through each square in board to initialize data attributes 
    for (var i=0; i < dim; i++) {
        for (var j=0; j < dim; j++) {
            boardSquare = ('<rect '
                    + 'class="boardSquare"' 
                    + 'x=' + j*sqLength + ' '
                    + 'y=' + i*sqLength + ' '
                    + 'data-i=' + i + ' '
                    + 'data-j=' + j + ' '
                    + 'width=' + sqLength + ' '
                    + 'height=' + sqLength + ' '
                    + 'style="fill:transparent; stroke-width:0.5; stroke:rgb(0,0,0)" />'
            );
            svg.innerHTML += boardSquare;
            // Render empty circles and initialize data attributes
            emptyPiece = ('<circle ' 
                + 'class="boardPiece" '
                + 'id="piece-' + i + j +'" ' 
                + 'data-value="0"'
                + 'data-i=' + i + ' '
                + 'data-j=' + j + ' '
                + 'cx='+ (j*sqLength + sqLength/2) +' ' 
                + 'cy=' + (i*sqLength + sqLength/2) +' '
                + 'r=' + (0.9*sqLength/2) +' ' 
                + 'fill=transparent '
                + '/>'
            );
            svg.innerHTML += emptyPiece;
        }
    }
}

function boardFreeze() {
    /** 
    * Disables click events on board
    */
    document.getElementById("board").style.pointerEvents = "none";
}

function boardThaw() {
    /** 
    * Enables click events on board
    */
    document.getElementById("board").style.pointerEvents = "";
}

function getPiece(i,j) {
    /** 
    * Gets specific circle svg piece with given coordinates on board
    * @param {number} i - ith row on board
    * @param {number} j - jth column on board
    * @returns {Object} circle svg for (i,j) position on board 
    */
    return document.getElementById("piece-"+i.toString()+j.toString());
}

function getNeighbour(i, j) {
    /** 
    * Given a coordinate, get positions of all immediate neighbouring empty positions
    * @param {number} i - ith row on board
    * @param {number} j - jth column on board
    * @returns {Array} Array of arrays containing empty neighbour coordinates
    */
    var svg = document.getElementById("board");
    // Initialize 8-axis directions: N-NE-E-SE-S-SW-W-NW
    var directions = [[-1,0],[-1,1],[0,1],[1,1],[1,0],[1,-1],[0,-1],[-1,-1]];
    var neighbours = []; // Initialize array to store neighbour coordinates

    for (var d = 0; d < directions.length ; d++) {
        var iNext  = parseInt(i) + directions[d][0];
        var jNext = parseInt(j) + directions[d][1];
        if (iNext > dim-1 || jNext > dim -1 || iNext < 0 || jNext < 0) {
            continue;
        } else {
            if (getPiece(iNext, jNext).dataset.value === "0") {
                neighbours.push(iNext.toString() + jNext.toString());
            }                           
        }
    }
    return neighbours;
}

function initNeighbours() {
    /** 
    * Initialize array of all empty spaces that are immediate neighbours of
    * all filled pieces on board
    * Updates data attribute in svg board to contain the array of empty neighbours
    */
    var svg = document.getElementById("board");
    var turn = parseInt(svg.dataset.turn); 
    var filled = svg.dataset.filled.split(","); 
    var allNeighbours = [];
    // Loop through all filled pieces and get their empty neighbours
    for (var k = 0; k < filled.length; k++) {
        iFilled = parseInt(filled[k].split("")[0]);
        jFilled = parseInt(filled[k].split("")[1]);
        if (parseInt(getPiece(iFilled, jFilled).dataset.value) === -turn) {
            neighbours = getNeighbour(iFilled, jFilled);
            for (var n = 0; n < neighbours.length; n++) {
                if (!inArray(neighbours[n], allNeighbours)) {
                    allNeighbours.push(neighbours[n]);
                }
            }    
        }
    }
    svg.dataset.neighbours = allNeighbours;
}

function checkValidDir(i, j, iDir, jDir) {
    /** 
    * Check in the direction (iDir, jDir), whether a valid move is available,
    * for a given empty space (i,j).
    * (i.e., the pieces will flip when move is made in that axis)
    * Direction is defined by (iDir, jDir), for example,
    *   if iDir=1, jDir=1), we are checking the N-E direction of a given space 
    *   if iDir=0, jDir=-1), we are checking the West direction of a given space 
    * @param {number} i - ith row on board
    * @param {number} j - jth column on board
    * @param {number} iDir - north/south direction
    * @param {number} jDir - east/west direction
    * @returns {boolean} true if the direction has a valid move to make
    */ 
    var valid = false;
    var svg = document.getElementById('board');
    var piece = getPiece(i, j);

    if (!parseInt(piece.dataset.value) === 0) {
        return false;
    }

    var firstShift = true; //Evaluating first neighbour in the iDir jDir direction

    while (true) {
        /* Starting form (i,j), traverse in the direction (iDir, jDir) and
        evaluate each piece to check if a valid move exists in the direction.*/
        i = parseInt(i) + iDir;
        j = parseInt(j) + jDir;

        var currPiece = getPiece(i, j);

        if (i > dim-1 || j > dim-1 || i < 0 || j < 0) {
            /* When coordinates are out of range, (i.e., out of the board),
            break the loop and return false for this direction's validity */
            break;
        } else if (currPiece.dataset.value === "0") {
            /* If position is empty, means no pieces to connect,
            then break out of the loop and return false for validity */ 
            break;
        } else if (currPiece.dataset.value === svg.dataset.turn && firstShift === true ) {
            /* If the currently evaluated space is same colour as the piece that is 
            about to be placed on (i,j), then this direction is invalid. 
            This only applies when current space is the first neighbour of the original starting space */
            break ; 
        } else if (currPiece.dataset.value === svg.dataset.turn && firstShift === false) {
            /* If the current space is same colour, but shifted before,
            means that previous space is different colour, end the check and this
            is where all pieces in this direction up to this piece is flipped.
            This direction will then be valid, and we return true */
            valid = true;
            break ;
        } else if (parseInt(currPiece.dataset.value) === -parseInt(svg.dataset.turn)) {
            /* If none of the above conditions are met, continue to traverse
            in the direction of (iDir, jDir) until condition is met to break
            the loop */
            firstShift = false;
            continue;
        }
    }
    return valid;
}

function checkValid(i, j) {
    /** 
    * Wrapper function that calls checkValidDir for 
    * every direction given a position (i,j) to be evaluated 
    * @param {number} i - ith row on board
    * @param {number} j - jth column on board
    * @returns {boolean} true if there exists at least one direction with a valid move
    */
    directions = [[-1,0],[-1,1],[0,1],[1,1],[1,0],[1,-1],[0,-1],[-1,-1]];
    for (var d=0; d < directions.length; d++) {
        if (checkValidDir(i,j, directions[d][0], directions[d][1])) {
            return true;
        }
    }
    return false;
}

function getAllValid() {
    /**
    * Evaluates every empty neighbour position on board
    * and calls checkValid on each position.
    * Collects all valid moves and stores it in data attribute
    * on svg board
    */
    var svg = document.getElementById("board");
    var allNeighbours = svg.dataset.neighbours.split(",");
    var validSpace = [];
    // Loops through each empty neighbour position and checks valid moves
    if (allNeighbours.length >= 1 && allNeighbours[0]!= "" ) {
        for (var n=0; n < allNeighbours.length; n++) {
            if(checkValid(parseInt(allNeighbours[n][0]), parseInt(allNeighbours[n][1]))) {
                validSpace.push(allNeighbours[n]);
            }
        }
    }
    svg.dataset.validspace = validSpace;
}

function renderBoard(board, turn) {
    /* 
    * Renders out all current black and white pieces on board
    * Updates scores on svg data attributes
    * Initializes all empty neighbours set
    * Initializes all valid moves available
    * @param {Array} board - 2D array of arrays representing board positions
    * @param {number} turn - Whose turn it is (1 or -1) / (black or white) 
    */
    //Set which side's turn 1 for black, -1 for white
    var svg = document.getElementById("board");
    svg.dataset.turn = turn;
    var blackscore = 0;
    var whitescore = 0;
    var filled = [];

    for (var i=0; i < board.length; i++) {
        for (var j=0; j < board.length; j++) {
            var pieceVal = board[i][j];
            var piece = getPiece(i,j);

            if (pieceVal === 1) {
                piece.style.fill = "black";
                piece.style.stroke ="black";
                piece.dataset.value = "1";
                filled.push(i.toString() + j.toString());
                blackscore += 1;
            }  
            if (pieceVal === -1) {
                piece.style.fill = "white";
                piece.style.stroke ="black";
                piece.dataset.value = "-1";
                filled.push(i.toString() + j.toString());
                whitescore += 1;
            }   
        }
    }
    svg.dataset.blackscore = blackscore;
    svg.dataset.whitescore = whitescore;
    svg.dataset.filled = filled;
    initNeighbours();
    getAllValid();
}

function Move(i,j, agent) {
    /** 
    * Places piece on board and flips opponent pieces 
    * Updates svg board and scores as well as valid moves
    * for the next user
    * @param {number} i - ith row on board
    * @param {number} j - jth column on board
    * @param {boolean} agent - true if move is made by agent 
    */
    agent = (typeof agent !== 'undefined') ? agent: false;

    var svg = document.getElementById("board");
    var turn = svg.dataset.turn;
    var filled = svg.dataset.filled.split(",");
    var validspace = svg.dataset.validspace.split(",");
    var directions = [[-1,0],[-1,1],[0,1],[1,1],[1,0],[1,-1],[0,-1],[-1,-1]] ;
    var flippedCount = 0;
    var piece = getPiece(i,j);  

    if (!inArray(i.toString()+j.toString(), validspace)) {
        // Do nothing if move is invalid.
        console.log("This is an invalid move");
    } else {
        // Flip pieces only if move is a valid one
        boardFreeze();
        if (turn === "1") {
            piece.style.fill = "black";
            piece.style.stroke ="black";
            piece.dataset.value = "1";
        }
        if (turn === "-1") {
            piece.style.fill = "white";
            piece.style.stroke ="black";
            piece.dataset.value = "-1";
        }
        // Update filled pieces array
        filled.push(i.toString() + j.toString());
        svg.dataset.filled = filled;

        for (var d=0; d < directions.length; d++) {
            // Flip all opposing pieces in the direction until same colour is met. 
            var iTemp = parseInt(i);
            var jTemp = parseInt(j);
            if (checkValidDir(i, j, directions[d][0], directions[d][1])) {
                while (true) {
                    iTemp += directions[d][0];
                    jTemp += directions[d][1];
                    nextPiece = getPiece(iTemp, jTemp);
                    if (parseInt(nextPiece.dataset.value) === parseInt(turn)) {
                        break;
                    } else {
                        //Flip the piece and move to next position
                        piece.dataset.value = turn;
                        nextPiece.dataset.value = turn;
                        if (turn === "1") {
                            piece.style.fill = "black";
                            nextPiece.style.fill = "black";
                            nextPiece.style.stroke ="black";
                        }
                        if (turn === "-1") {
                            piece.style.fill = "white";
                            nextPiece.style.fill = "white";
                            nextPiece.style.stroke ="black";
                        }
                        flippedCount += 1;
                    }
                }
            }
        }
        // Update scores after all pieces flipped
        // Total score increase = all flipped pieces + 1 new piece placed
        if (turn === "1") {
            svg.dataset.blackscore = parseInt(svg.dataset.blackscore) + flippedCount + 1;
            svg.dataset.whitescore = parseInt(svg.dataset.whitescore) - flippedCount; 
        }   
        if (turn === "-1") {
            svg.dataset.whitescore = parseInt(svg.dataset.whitescore) + flippedCount + 1;
            svg.dataset.blackscore = parseInt(svg.dataset.blackscore) - flippedCount;
        }
        document.getElementById("blackscore").innerHTML = svg.dataset.blackscore;
        document.getElementById("whitescore").innerHTML = svg.dataset.whitescore;
        
        turn = -parseInt(turn); // Flip to next turn
        svg.dataset.turn = turn;
        initNeighbours(); // Update latest neighbour pieces
        getAllValid(); // Update latest valid positions

        /* The following code below evaluetes the scenario where there are no
        valid moves for the next player. The next player will then skip their turn. 
        When both players have no more moves to make, the game ends.*/ 

        validspace = svg.dataset.validspace.split(","); // Get latest valid spaces

        // If there are no valid moves for the player, skip turn 
        if (validspace.length === 1 && validspace[0] == "") {
            turn = -parseInt(turn);
            svg.dataset.turn = turn;
            initNeighbours();
            getAllValid();
            validspace = svg.dataset.validspace.split(",");
            // If both players have no more moves
            if (validspace.length === 1 && validspace[0] == "") {
                blackscoreFinal = parseInt(svg.dataset.blackscore);
                whitescoreFinal = parseInt(svg.dataset.whitescore);
                if (blackscoreFinal > whitescoreFinal) {
                    alert("Black wins: ", 
                        blackscoreFinal.toString(), 
                        " : ", 
                        whitescoreFinal.toString()
                    );
                } else if (blackscoreFinal < whitescoreFinal) {
                    alert(
                        "White wins: ", blackscoreFinal.toString(), 
                        " : ", 
                        whitescoreFinal.toString()
                    );
                } else {
                    alert(" Draw: ", 
                        blackscoreFinal.toString(), " : ", 
                        whitescoreFinal.toString()
                    );
                }
            } else {
                if (agent) {
                    // If the move was made by the agent, the agent gets to moves again
                    updateMessage("No move for you! Agent's turn");
                    agentMove(turn);
                } else {
                    // If the move was made by user, user gets to move again
                    updateMessage("No moves for agent, your turn again");
                    boardThaw(); //Unfreeze board for user's turn again
                }
            }
        } else {
            if (!agent) {
                // If the move was not made by the agent,
                // It is agent's turn to move
                agentMove(turn);
            }
        }
    }
}

function getColourArrays() {
    /** 
    * Helper function to obtain arrays of black and white positions
    * from the current game board
    * Used for consumption of agentMove fuction API endpoint
    * @example
    * // example output
    * {
    *   blackArray: [[3,3], [3,4], [4,3], [2,3]],
    *   whiteArray: [[4,4], [2,4], [3,5], [2,6]]
    * }
    * @returns {Object} two arrays containing positions of black/white pieces 
    */      
    var filled = document.getElementById("board").dataset.filled.split(",");
    var blackArray = [];
    var whiteArray = [];
    // Iterate through filled pieces and sort positions into 
    // blackArray and whiteArray
    for (k=0; k < filled.length; k++) {
        var piece = document.getElementById("piece-"+filled[k]);
        if (piece.dataset.value == 1) {
            blackArray.push([parseInt(piece.dataset.i), parseInt(piece.dataset.j)]);
        }
        if (piece.dataset.value == -1) {
            whiteArray.push([parseInt(piece.dataset.i), parseInt(piece.dataset.j)]);
        }
    }
    return {
        "blackArray":blackArray,    
        "whiteArray":whiteArray
    }
}

function agentMove(turn) {
    /**
    * Makes AJAX POST request to application obtain agent's next Move 
    * @param {number} turn - Integer 1 or -1 representing the Agent's colour
    */
    updateMessage("Thinking...");
    // Prepare black and white piece positions for POST request payload
    var colourArrays = getColourArrays();
    // Create an XMLHttpRequest for API call
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState === 4 && this.status === 200) {
            // Parse response from backend and make the move on the svg board
            var json_response = JSON.parse(this.response);
            Move(json_response['move'][0], json_response['move'][1], agent=true);
            boardThaw(); // Always unfreeze board for user to play next
            updateMessage("Your Turn ("+turn2Colour(-turn)+")")
            return json_response
        }
    }
    xhttp.open('POST', '/search_move', true);
    xhttp.setRequestHeader('Content-Type', 'application/json');
    xhttp.send(
    JSON.stringify({
        "blackFilled":colourArrays["blackArray"],    
        "whiteFilled":colourArrays["whiteArray"],    
        "turn":turn,
        })
    );
}

width =800; 
dim = 8;
drawGrid(800, 8);
boardPieces = document.getElementsByClassName("boardPiece")
boardSquares = document.getElementsByClassName("boardSquare")

for (m=0; m < boardSquares.length; m++ ){

    boardSquares[m].addEventListener("click", function() {
        Move(this.dataset.i, this.dataset.j);
    })
}
for (m=0; m < boardPieces.length; m++ ){

    boardPieces[m].addEventListener("click", function() {
        Move(this.dataset.i, this.dataset.j);
    })
}
// document.getElementsByTagName('')ElementById("board").innerHTML='<circle cx="100" cy="200" r="40" stroke="black" stroke-width="4" fill="yellow" />'
b = newBoard(dim);
renderBoard(b, "1");
// boardFreeze();
// agentMove(1);
initNeighbours();