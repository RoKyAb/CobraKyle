package experiment

import (
	"battlesnake/appengine/game"
	"math"
	"math/rand"
	"sort"
)

func foodBonuses(moves []string, b game.Battlesnake, foods []game.Coord, w game.Weights) map[string]float64 {
	type foodDistancePair struct {
		Move     string
		Distance float64
	}

	foodDistance := make([]foodDistancePair, len(moves))
	for i, m := range moves {
		newHead := movedHead(b.Head, m)
		closestFoodThisWay := 1000.0
		for _, f := range foods {
			d := lineDistance(f, newHead)
			if d < closestFoodThisWay {
				closestFoodThisWay = d
			}
		}
		foodDistance[i] = foodDistancePair{Move: m, Distance: closestFoodThisWay}
	}

	sort.Slice(foodDistance, func(i, j int) bool {
		return foodDistance[i].Distance < foodDistance[j].Distance
	})

	result := make(map[string]float64)

	closestFood := foodDistance[0].Distance
	for i, fd := range foodDistance {
		if i == 0 {
			result[fd.Move] = w.FoodBonus1
		}

		if i == 1 {
			result[fd.Move] = w.FoodBonus2
		}

		if i == 2 {
			result[fd.Move] = w.FoodBonus3
		}

		if fd.Distance == closestFood {
			result[fd.Move] = w.FoodBonus1
		}
	}

	return result
}

func lowRiskMove(moves []string, me game.Battlesnake, board game.Board, w game.Weights) string {
	foodBonusMap := foodBonuses(moves, me, board.Food, w)

	points := math.MaxFloat64
	move := ""
	for _, m := range moves {
		newHead := movedHead(me.Head, m)
		newBody := append([]game.Coord{newHead}, me.Body...)
		var nPoints float64
		nPoints += bodyTightnessBonus(newHead, newBody, w)
		nPoints += adjacentToWall(newHead, m, board.Height, board.Width, w)
		nPoints += opponentProximity(me.ID, newHead, me.Length, board.Snakes, w)

		if me.Health < int32(len(board.Snakes)*10) {
			nPoints += foodBonusMap[m]
		}

		if me.Health < int32(len(board.Snakes)*5) {
			nPoints += foodBonusMap[m]
		}

		//fmt.Println(fmt.Sprintf("%s: %d [s:%d w:%d o:%d] (%s)", m, nPoints, self, wallD, op, me.ID))
		//fmt.Println(fmt.Sprintf("%s: %d (%s)", m, nPoints, me.ID))

		if nPoints < points {
			points = nPoints
			move = m
		}

		if nPoints == points && rand.Float64() > 0.5 {
			move = m
		}
	}
	return move
}

func bodyTightnessBonus(head game.Coord, body []game.Coord, w game.Weights) float64 {
	var bonus float64
	for _, b := range body {
		if nearby(head, b) {
			bonus += w.TightSnakeBonus
		}
		if district(head, b) {
			bonus += w.TightSnakeBonus/2
		}
	}

	return bonus
}

func opponentProximity(myID string, head game.Coord, myLength int32, snakes []game.Battlesnake, weights game.Weights) float64 {
	var points float64
	for _, s := range snakes {
		if s.ID != myID {
			if s.Length > myLength {
				if nearby(head, s.Head) {
					points += weights.OpponentHeadPenalty
				}

				if district(head, s.Head) {
					points += weights.OpponentHeadPenalty / 2
				}
			}
			for _, b := range s.Body {
				if nearby(head, b) {
					points += weights.OpponentBodyPenalty
				}

				if district(head, b) {
					points += weights.OpponentBodyPenalty/2
				}
			}
		}
	}

	return points
}

func lineDistance(a, b game.Coord) float64 {
	return math.Sqrt(math.Pow(math.Abs(float64(a.X-b.X)), 2) + math.Pow(math.Abs(float64(a.Y-b.Y)), 2))
}

func adjacent(head, bodyPart game.Coord) bool {
	return 1 == lineDistance(head, bodyPart)
}

func nearby(head, bodyPart game.Coord) bool {
	return math.Sqrt2 >= lineDistance(head, bodyPart)
}

func district(head, bodyPart game.Coord) bool {
	return 2 >= lineDistance(head, bodyPart)
}

func adjacentToWall(head game.Coord, move string, h int, w int, weights game.Weights) float64 {
	if move == "up" && head.Y >= h-1 {
		return weights.WallPenalty
	}
	if move == "down" && head.Y <= 0 {
		return weights.WallPenalty
	}
	if move == "left" && head.X <= 0 {
		return weights.WallPenalty
	}
	if move == "right" && head.X >= w-1 {
		return weights.WallPenalty
	}
	return 0
}

func dontHitWallOrSelfOrOpponents(b game.Battlesnake, board game.Board) []string {
	possibleMoves := []string{}
	allMoves := []string{"up", "down", "left", "right"}

	// Try avoid everything
	for _, m := range allMoves {
		newHead := movedHead(b.Head, m)
		if missWalls(newHead, board.Height, board.Width) && missSelf(newHead, b.Body) && missOpponents(b.ID, newHead, board.Snakes) && dontEnclose(newHead, b, board) {
			possibleMoves = append(possibleMoves, m)
		}
	}

	// At least try not to hit on next
	if len(possibleMoves) == 0 {
		for _, m := range allMoves {
			newHead := movedHead(b.Head, m)
			if missWalls(newHead, board.Height, board.Width) && missSelf(newHead, b.Body) && missOpponents(b.ID, newHead, board.Snakes) {
				possibleMoves = append(possibleMoves, m)
			}
		}
	}

	// At least try not to hit wall or self
	if len(possibleMoves) == 0 {
		for _, m := range allMoves {
			newHead := movedHead(b.Head, m)
			if missWalls(newHead, board.Height, board.Width) && missSelf(newHead, b.Body) {
				possibleMoves = append(possibleMoves, m)
			}
		}
	}

	if len(possibleMoves) == 0 {
		return allMoves
	}

	return possibleMoves
}

func missWalls(newHead game.Coord, h int, w int) bool {
	return newHead.X >= 0 && newHead.X < w && newHead.Y >= 0 && newHead.Y < h
}

func missSelf(newHead game.Coord, body []game.Coord) bool {
	for _, b := range body {
		if b == newHead {
			return false
		}
	}
	return true
}

func missOpponents(myID string, newHead game.Coord, snakes []game.Battlesnake) bool {
	for _, s := range snakes {
		if s.ID == myID {
			continue
		}
		for _, b := range s.Body {
			if b == newHead {
				return false
			}
		}
	}
	return true
}

func dontEnclose(newHead game.Coord, b game.Battlesnake, board game.Board) bool {
	possibleMoves := []string{}
	newBody := append([]game.Coord{newHead}, b.Body[:len(b.Body)-1]...)
	allMoves := []string{"up", "down", "left", "right"}
	for _, m := range allMoves {
		futureHead := movedHead(newHead, m)
		if missWalls(futureHead, board.Height, board.Width) && missSelf(futureHead, newBody) && missOpponents(b.ID, futureHead, board.Snakes) {
			possibleMoves = append(possibleMoves, m)
		}
	}

	if len(possibleMoves) == 0 {
		return false
	}

	possibleMoves2 := []string{}
	for _, m := range possibleMoves {
		pm := []string{}
		futureHead := movedHead(newHead, m)
		newNewBody := append([]game.Coord{futureHead}, newBody[:len(newBody)-1]...)
		for _, m := range allMoves {
			futureFutreHead := movedHead(futureHead, m)
			if missWalls(futureFutreHead, board.Height, board.Width) && missSelf(futureFutreHead, newNewBody) && missOpponents(b.ID, futureFutreHead, board.Snakes) {
				pm = append(pm, m)
			}
		}
		if len(pm) > len(possibleMoves2) {
			possibleMoves2 = pm
		}
	}

	return len(possibleMoves2) > 0
}

func movedHead(head game.Coord, move string) game.Coord {
	delta := map[string]game.Coord{
		"up":    {X: 0, Y: 1},
		"down":  {X: 0, Y: -1},
		"left":  {X: -1, Y: 0},
		"right": {X: 1, Y: 0},
	}
	return game.Coord{
		X: head.X + delta[move].X,
		Y: head.Y + delta[move].Y,
	}
}
