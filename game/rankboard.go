package game

import (
	"fmt"
	"strings"
)

//TODO
// Implement rank over score
// Rank is a map of player name to the rankMap score
// Rank is calculated by how many people a play bests in all previous games:
// Winner of a game of 5 will get a rankMap score of 4 where last place will get a rankMap score of 0
var rankMap map[string]int

func InitializeRank() {
	rankMap = make(map[string]int)
	rankMap["snowie"] = 98
	rankMap["xuechenf"] = 99
}

func PrintRank() string{
	builder := strings.Builder{}

	for k, v := range rankMap {
		builder.WriteString(fmt.Sprintf("%s -> %d <br />\n\r", k, v)) }
	return builder.String()
}

func UpdatePlayerRanks(players []*Player) {
	for _, p := range players {
		rankMap[p.Name] += p.Score // doesn't gives 0 as value and an error, toss away error and keep 0 \o/
	}
}