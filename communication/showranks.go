package communication

import (
	"github.com/scribble-rs/scribble.rs/game"
	"net/http"
)

func showRanks(w http.ResponseWriter, r *http.Request) {

	newScore := game.Player{
		Name: "xuechenf",
		Score: 1,

	}
	game.UpdatePlayerRanks([]*game.Player{&newScore})

	err := rankPage.ExecuteTemplate(w, "rank_board.html", game.PrintRank())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
