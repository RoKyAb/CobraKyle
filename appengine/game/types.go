package game

import "github.com/BattlesnakeOfficial/rules"

type Game struct {
	ID      string `json:"id"`
	Timeout int32  `json:"timeout"`
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Battlesnake struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Health int32   `json:"health"`
	Body   []Coord `json:"body"`
	Head   Coord   `json:"head"`
	Length int32   `json:"length"`
	Shout  string  `json:"shout"`
}

type Board struct {
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Food   []Coord       `json:"food"`
	Snakes []Battlesnake `json:"snakes"`
}

type BattlesnakeInfoResponse struct {
	APIVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
}

type GameRequest struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

type MoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout,omitempty"`
}

type Weights struct {
	FoodBonus1 				float64
	FoodBonus2 				float64
	FoodBonus3 				float64
	TightSnakeBonus 		float64
	WallPenalty 			float64
	OpponentHeadPenalty 	float64
	OpponentBodyPenalty 	float64
}

func PointToCoord(point rules.Point) Coord {
	return Coord{
		X: int(point.X),
		Y: int(point.Y),
	}
}

func RuleSnekToBattleSnek(snake rules.Snake) Battlesnake {
	body := make([]Coord, len(snake.Body))
	for i, p := range snake.Body {
		body[i] = PointToCoord(p)
	}

	return Battlesnake{
		ID:     snake.ID,
		Name:   snake.ID,
		Health: snake.Health,
		Body:   body,
		Head:   PointToCoord(snake.Body[0]),
		Length: int32(len(snake.Body)),
		Shout:  "",
	}
}

func RuleBoardStateToGameReq(state *rules.BoardState) GameRequest {
	food := make([]Coord, len(state.Food))
	for i, p := range state.Food {
		food[i] = PointToCoord(p)
	}

	snakes := make([]Battlesnake, len(state.Snakes))
	for i, s := range state.Snakes {
		snakes[i] = RuleSnekToBattleSnek(s)
	}

	return GameRequest{
		Board: Board{
			Height: int(state.Height),
			Width:  int(state.Width),
			Food:   food,
			Snakes: snakes,
		},
		You:   Battlesnake{},
	}
}