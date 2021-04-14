package image

import (
	"SecretHitlerDiscordgo/Game"
	"github.com/bwmarrin/discordgo"
	"github.com/nakabonne/gosivy/agent"
	"os"
	"testing"
	"time"
)

func init() {
	_ = os.Mkdir("./temp", 0755)
	_ = os.Mkdir("./avatars", 0755)

	if err := agent.Listen(agent.Options{}); err != nil {
		panic(err)
	}
	defer agent.Close()

	time.Sleep(2 * time.Second)
}

func TestDrawFascistBoard(t *testing.T) {
	g := Game.NewGame()
	g.FascistTracker = 6

	DrawFascistBoard(g).SavePNG("./temp/fascist-3.png")
}

func TestDrawLiberalBoard(t *testing.T) {
	g := Game.NewGame()
	g.ElectionTracker = 3
	g.LiberalTracker = 5

	DrawLiberalBoard(g).SavePNG("./temp/liberal-3-2.png")
}

func TestDrawStatus(t *testing.T) {
	var (
		users = []discordgo.User{
			{
				ID:       "145618075452964864",
				Username: "TheTipo01",
				Avatar:   "93d255afb6f8d89fab55360edad0a9ef",
			},
			{
				ID:       "143060848091463680",
				Username: "dany_ev3",
				Avatar:   "a0527abc7a7a3674529c6271bcc15f16",
			},
			{
				ID:       "409711680633700373",
				Username: "techmccat",
				Avatar:   "13518517e19c32a9bf8e2cc740c2015e",
			},
			{
				ID:       "148395955962511360",
				Username: "Hexa",
				Avatar:   "f5452a1008bf89035c1661ba748a94f8",
			},
			{
				ID:       "322756205024116739",
				Username: "\U0001F9FF👄\U0001F9FF",
				Avatar:   "0fd0609328558c514d4edc2574f79691",
			},
			{
				ID:       "783071008164282439",
				Username: "Michele Bolla",
				Avatar:   "0656e7420082e5adaabe2afe7afb4244",
			},
			{
				ID:       "271001798473416704",
				Username: "Xx_DNS_xX",
				Avatar:   "e77e1dbc885595545f47c74bdda6dec0",
			},
			{
				ID:       "145874051032678400",
				Username: "slashtube",
				Avatar:   "ee2bc862adc078bd5814ba4bbb2d96f5",
			},
		}

		g = Game.NewGame()
	)

	g.Players = make([]Game.Player, 8)

	for i, u := range users {
		DownloadAvatar(&u)
		g.Players[i] = Game.Player{
			Id:   u.ID,
			Role: 0,
		}
	}

	g.Players[0].Role = Game.FascistRole
	g.Players[1].Role = Game.FascistRole
	g.Players[2].Role = Game.HitlerRole
	g.Players[3].Role = Game.LiberalRole
	g.Players[4].Role = Game.LiberalRole
	g.Players[5].Role = Game.FascistRole
	g.Players[6].Role = Game.LiberalRole
	g.Players[7].Role = Game.FascistRole

	g.Chancellor = &g.Players[0]
	g.President = &g.Players[5]

	DrawStatus(g, &g.Players[2]).SavePNG("./temp/statusHitler.png")
	DrawStatus(g, &g.Players[0]).SavePNG("./temp/statusFascist.png")
	DrawStatus(g, &g.Players[3]).SavePNG("./temp/statusLiberal.png")
}