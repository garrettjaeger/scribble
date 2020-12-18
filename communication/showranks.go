package communication

import (
	"github.com/scribble-rs/scribble.rs/game"
	"net/http"
)

func showRanks(w http.ResponseWriter, r *http.Request) {

	//newScore := game.Player{
	//	Name: "xuechenf",
	//	Score: 1,
	//
	//}
	//fmt.Println("Updating score")
	//fmt.Println(game.UpdatePlayerScores([]*game.Player{&newScore}, nodeid.NodeID))
	//fmt.Println("done")

	err := rankPage.ExecuteTemplate(w, "rank_board.html", game.PrintScore())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
