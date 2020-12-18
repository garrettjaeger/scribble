package communication

import (
	"fmt"
	"github.com/scribble-rs/scribble.rs/game"
	"github.com/scribble-rs/scribble.rs/nodeid"
	"html/template"
	"net/http"
	"strings"
)

func increaseScore(w http.ResponseWriter, r *http.Request) {

	var newScores []*game.Player
	newScores = append(newScores, &game.Player{
		Name: "meowlady",
		Score: 2,
	})
	newScores = append(newScores, &game.Player{
		Name: "woofman",
		Score: 3,
	})
	newScores = append(newScores, &game.Player{
		Name: "garjaeg",
		Score: 4,
	})
	newScores = append(newScores, &game.Player{
		Name: "xuechenf",
		Score: 1,
	})

	fmt.Println("Updating score")
	fmt.Println(game.UpdatePlayerScores(newScores, nodeid.NodeID))
	fmt.Println("done")

	builder := strings.Builder{}
	builder.WriteString("New scored added: <br />\n\r")
	for _, v := range newScores {
		builder.WriteString(fmt.Sprintf("%s -> %d <br />\n\r", v.Name, v.Score)) }

	err := scoreIncrementPage.ExecuteTemplate(w, "increments_score.html", template.HTML(builder.String()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
