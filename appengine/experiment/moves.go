package experiment

import (
	"math"
	"math/rand"
	"sort"
)

func foodBonuses(moves []string, b Battlesnake, foods []Coord) map[string]int {
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

	result := make(map[string]int)

	closestFood := foodDistance[0].Distance
	for i, fd := range foodDistance {
		if fd.Distance == closestFood {
			result[fd.Move] = -6
		}
		result[fd.Move] = -5 + i
	}

	return result
}

func lowRiskMove(moves []string, me Battlesnake, board Board) string {
	foodBonusMap := foodBonuses(moves, me, board.Food)

	points := 1000
	move := ""
	for _, m := range moves {
		newHead := movedHead(me.Head, m)
		newBody := append([]Coord{newHead}, me.Body...)
		nPoints := 0
		for _, b := range newBody {
			if !nearby(newHead, b) {
				nPoints += 1
			}
		}

		nPoints += adjacentToWall(newHead, m, board.Height, board.Width)
		nPoints += opponentProximity(me.ID, newHead, me.Length, board.Snakes)

		if me.Health < int32(len(board.Snakes)*20) {
			nPoints += foodBonusMap[m]
		}

		if me.Health < int32(len(board.Snakes)*10) {
			nPoints += foodBonusMap[m]
		}

		//fmt.Println(fmt.Sprintf("%s: %d [s:%d w:%d o:%d] (%s)", m, nPoints, self, wallD, op, me.ID))
		//fmt.Println(fmt.Sprintf("%s: %d (%s)", m, nPoints, me.ID))

		if nPoints < points {
			points = nPoints
			move = m
		}

		if nPoints == points && rand.Float64() > 0.5 {
			points = nPoints
			move = m
		}
	}
	return move
}

func opponentProximity(myID string, head Coord, myLength int32, snakes []Battlesnake) int {
	points := 0
	for _, s := range snakes {
		if s.ID != myID {
			if withInThree(head, s.Head) && s.Length > myLength {
				points += 5
			}

			for _, b := range s.Body {
				if nearby(head, b) {
					points += 1
				}
			}
		}
	}

	return points
}

func lineDistance(a, b Coord) float64 {
	return math.Sqrt(math.Pow(math.Abs(float64(a.X-b.X)), 2) + math.Pow(math.Abs(float64(a.Y-b.Y)), 2))
}

func adjacent(head, bodyPart Coord) bool {
	return 1 == lineDistance(head, bodyPart)
}

func nearby(head, bodyPart Coord) bool {
	return math.Sqrt2 >= lineDistance(head, bodyPart)
}

func withInThree(head, bodyPart Coord) bool {
	return math.Sqrt(3) >= lineDistance(head, bodyPart)
}

func adjacentToWall(head Coord, move string, h int, w int) int {
	if move == "up" && head.Y >= h-1 {
		return 3
	}
	if move == "down" && head.Y <= 0 {
		return 3
	}
	if move == "left" && head.X <= 0 {
		return 3
	}
	if move == "right" && head.X >= w-1 {
		return 3
	}
	return 0
}

func dontHitWallOrSelfOrOpponents(b Battlesnake, board Board) []string {
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

func missWalls(newHead Coord, h int, w int) bool {
	return newHead.X >= 0 && newHead.X < w && newHead.Y >= 0 && newHead.Y < h
}

func missSelf(newHead Coord, body []Coord) bool {
	for _, b := range body {
		if b == newHead {
			return false
		}
	}
	return true
}

func missOpponents(myID string, newHead Coord, snakes []Battlesnake) bool {
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

func dontEnclose(newHead Coord, b Battlesnake, board Board) bool {
	possibleMoves := []string{}
	newBody := append([]Coord{newHead}, b.Body[:len(b.Body)-1]...)
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
		newNewBody := append([]Coord{futureHead}, newBody[:len(newBody)-1]...)
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

func movedHead(head Coord, move string) Coord {
	delta := map[string]Coord{
		"up":    {X: 0, Y: 1},
		"down":  {X: 0, Y: -1},
		"left":  {X: -1, Y: 0},
		"right": {X: 1, Y: 0},
	}
	return Coord{
		X: head.X + delta[move].X,
		Y: head.Y + delta[move].Y,
	}
}