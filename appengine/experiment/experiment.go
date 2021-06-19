package experiment

import (
	"battlesnake/appengine/game"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HandleIndex is called when your Battlesnake is created and refreshed
// by play.battlesnake.com. BattlesnakeInfoResponse contains information about
// your Battlesnake, including what it should look like on the game board.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := game.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "kabra",
		Color:      "#00FF00",
		Head:       "silly",
		Tail:       "hook",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

// HandleStart is called at the start of each game your Battlesnake is playing.
// The GameRequest object contains information about the game that's about to start.
// TODO: Use this function to decide how your Battlesnake is going to look on the board.
func HandleStart(w http.ResponseWriter, r *http.Request) {
	request := game.GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Print("START\n")
}

func MovesTest(g game.GameRequest, weights game.Weights) string {
	possibleMoves := dontHitWallOrSelfOrOpponents(g.You, g.Board)
	return lowRiskMove(possibleMoves, g.You, g.Board, weights)
}

// HandleMove is called for each turn of each game.
// Valid responses are "up", "down", "left", or "right".
// TODO: Use the information in the GameRequest object to determine your next move.
func HandleMove(w http.ResponseWriter, r *http.Request) {
	request := game.GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	possibleMoves := dontHitWallOrSelfOrOpponents(request.You, request.Board)

	ws := game.Weights{
		FoodBonus1:          -6,
		FoodBonus2:          -2,
		FoodBonus3:          -3,
		TightSnakeBonus:     -1,
		WallPenalty:          3,
		OpponentHeadPenalty:  3,
		OpponentBodyPenalty:  2,
	}

	move := lowRiskMove(possibleMoves, request.You, request.Board, ws)

	response := game.MoveResponse{
		Move: move,
	}

	//fmt.Printf("MOVE: %s\n", response.Move)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

// HandleEnd is called when a game your Battlesnake was playing has ended.
// It's purely for informational purposes, no response required.
func HandleEnd(w http.ResponseWriter, r *http.Request) {
	request := game.GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Print("END\n")
}
