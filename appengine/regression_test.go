package main

import (
	"battlesnake/appengine/game"
	"fmt"
	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/require"

	"battlesnake/appengine/experiment"
	"testing"
)

func TestWeights(t *testing.T) {

	ruleset := rules.Ruleset(&rules.StandardRuleset{})

	wSet := []game.Weights{
		{
			FoodBonus1:          -6,
			FoodBonus2:          -2,
			FoodBonus3:          -3,
			TightSnakeBonus:     -1,
			WallPenalty:          3,
			OpponentHeadPenalty:  3,
			OpponentBodyPenalty:  2,
		},
	}

	/*
	for _, ohp := range []float64{0, 1, 2, 3, 4, 5, 6, 7} {
		for _, obp := range []float64{0, 1, 2, 3, 4, 5, 6, 7} {

			wSet = append(wSet, game.Weights{
				FoodBonus1:          -6,
				FoodBonus2:          -2,
				FoodBonus3:          -3,
				TightSnakeBonus:     -1,
				WallPenalty:         3,
				OpponentHeadPenalty: ohp,
				OpponentBodyPenalty: obp,
			})
		}
	}
	 */


	fmt.Println(len(wSet))

	ids := []string{"Control", "Control2", "Experiment", "Experiment2"}

	type result struct {
		ws game.Weights
		ew int
		cw int
	}

	top20WSets := make([]result, 20)

	maybeAddWS := func(weights game.Weights, ew int, cw int) bool {
		for i , res := range top20WSets {
			if res.ew < ew {
				top20WSets[i] = result{ws: weights, ew: ew, cw: cw}
				return true
			}
		}
		return false
	}

	for _, ws := range wSet {
		ControlWins := 0
		ExperimentWins := 0

		for i := 0; i < 1000; i++ {
			bs, err := ruleset.CreateInitialBoardState(11, 11, ids)
			require.NoError(t, err)

			var gameOver bool
			j := 0
			winner := ""
			for !gameOver {

				nextMoves := []rules.SnakeMove{}

				food := make([]game.Coord, len(bs.Food))
				for i, p := range bs.Food {
					food[i] = pointToCoord(p)
				}

				snakes := make([]game.Battlesnake, len(bs.Snakes))
				for i, s := range bs.Snakes {
					snakes[i] = ruleSnekToBattleSnek(s)
				}

				for i, s := range bs.Snakes {
					req := game.GameRequest{
						Board: game.Board{
							Height: 11,
							Width:  11,
							Food:   food,
							Snakes: snakes,
						},
						You:   snakes[i],
					}

					if (s.ID == ids[0] || s.ID == ids[1]) && s.EliminatedCause == rules.NotEliminated {

						nextMoves = append(nextMoves, rules.SnakeMove{
							ID:   s.ID,
							Move: movesTest(req),
						} )
					}

					if (s.ID == ids[2] || s.ID == ids[3]) && s.EliminatedCause == rules.NotEliminated {

						nextMoves = append(nextMoves, rules.SnakeMove{
							ID:   s.ID,
							Move: experiment.MovesTest(req, ws),
						} )
					}
				}

				newBS, err := ruleset.CreateNextBoardState(bs, nextMoves)
				require.NoError(t, err)

				gameOver, err = ruleset.IsGameOver(newBS)
				require.NoError(t, err)

				if gameOver {
					for _, s := range newBS.Snakes {
						if (s.ID == ids[0] || s.ID == ids[1]) && s.EliminatedCause == rules.NotEliminated {
							ControlWins += 1
							winner = s.ID
						}

						if (s.ID == ids[2] || s.ID == ids[3]) && s.EliminatedCause == rules.NotEliminated {
							ExperimentWins += 1
							winner = s.ID
						}
					}
				}

				bs = newBS
				j += 1
			}
			fmt.Println(fmt.Sprintf("Game: %d Turns: %d Winner:%s", i, j, winner))
		}
		maybeAddWS(ws, ExperimentWins, ControlWins)
	}
	fmt.Println(top20WSets)
}

func pointToCoord(point rules.Point) game.Coord {
	return game.Coord{
		X: int(point.X),
		Y: int(point.Y),
	}
}

func ruleSnekToBattleSnek(snake rules.Snake) game.Battlesnake {
	body := make([]game.Coord, len(snake.Body))
	for i, p := range snake.Body {
		body[i] = pointToCoord(p)
	}

	return game.Battlesnake{
		ID:     snake.ID,
		Name:   snake.ID,
		Health: snake.Health,
		Body:   body,
		Head:   pointToCoord(snake.Body[0]),
		Length: int32(len(snake.Body)),
		Shout:  "",
	}
}