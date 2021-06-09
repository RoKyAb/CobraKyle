package main

import (
	"math"
	"math/rand"
)

func findFood(moves []string, b Battlesnake, foods []Coord ) string {
	closestFood := 1000.0
	move := ""
	for _, m := range moves {
		newHead := movedHead(b.Head, m)
		for _, f := range foods {
			d := lineDistance(f, newHead)
			if d < closestFood {
				move = m
				closestFood = d
			}
		}
	}
	return move
}

func lowRiskMove(moves []string, b Battlesnake, foods []Coord, h int, w int) string {
	if b.Health < 15 {
		return findFood(moves, b, foods)
	}

	density := 1000
	move := ""
	for _, m := range moves {
		newHead := movedHead(b.Head, m)
		newBody := append([]Coord{newHead}, b.Body...)
		nDensity := 0
		for _, b := range newBody {
			if adjacent(newHead, b) {
				nDensity += 1
			}
			nDensity += adjacentToWall(newHead, m, h, w)
		}
		if nDensity < density {
			density = nDensity
			move = m
		}

		if nDensity == density && rand.Float64() > 0.5 {
			density = nDensity
			move = m
		}
	}
	return move
}

func lineDistance(a, b Coord) float64 {
	return math.Sqrt(math.Pow(math.Abs(float64(a.X - b.X)),2) + math.Pow(math.Abs(float64(a.Y - b.Y)),2))
}

func adjacent(head, bodyPart Coord) bool {
	//return 1 <= lineDistance(head, bodyPart)
	return math.Sqrt2 <= lineDistance(head, bodyPart)
}

func adjacentToWall(head Coord, move string, h int, w int) int {
	if move == "up" && head.Y >= h-1 {
		return 3
	}
	if move == "down" && head.Y <= 0 {
		return 3
	}
	if move == "left" && head.X <= 0{
		return 3
	}
	if move == "right" && head.X >= w-1 {
		return 3
	}
	return 0
}

func dontHitWallOrSelf(b Battlesnake, h int, w int) []string {
	possibleMoves := []string{}
	allMoves := []string{"up", "down", "left", "right"}
	for _, m := range allMoves {
		newHead := movedHead(b.Head, m)
		if missWalls(newHead, h, w) && missSelf(newHead, b.Body) && dontEnclose(newHead, b.Body, w, h) {
			//fmt.Print(newHead)
			//fmt.Println(" "+m)
			possibleMoves = append(possibleMoves, m)
		}
	}

	//fmt.Print("Body: ")
	//fmt.Println(b.Body)

	//fmt.Print("Possible Moves: ")
	//fmt.Println(possibleMoves)

	// At least try not to hit wall or self
	if len(possibleMoves) == 0 {
		for _, m := range allMoves {
			newHead := movedHead(b.Head, m)
			if missWalls(newHead, h, w) && missSelf(newHead, b.Body) {
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

func dontEnclose(newHead Coord, body []Coord, h, w int) bool {
	possibleMoves := []string{}
	newBody := append([]Coord{newHead}, body[:len(body)-1]...)
	allMoves := []string{"up", "down", "left", "right"}
	for _, m := range allMoves {
		futureHead := movedHead(newHead, m)
		if missWalls(futureHead, h, w) && missSelf(futureHead, newBody) {
			//fmt.Print(newHead)
			//fmt.Println(" "+m)
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
			if missWalls(futureFutreHead, h, w) && missSelf(futureFutreHead, newNewBody) {
				//fmt.Print(newHead)
				//fmt.Println(" "+m)
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
		"up" : {X: 0, Y: 1},
		"down" : {X: 0, Y: -1},
		"left" : {X: -1, Y: 0},
		"right" : {X: 1, Y: 0},
	}
	return Coord{
		X: head.X + delta[move].X,
		Y: head.Y + delta[move].Y,
	}
}
