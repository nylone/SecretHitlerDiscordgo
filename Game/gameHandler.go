package Game

import (
	"SecretHitlerDiscordgo/Utils"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// collection of functions that handle the flow of the game

// NewGame creates the basic structure for the game to let others join using the Join command
func NewGame() *Game {
	G := Game{
		out:  make(chan interface{}),
		in:   make(chan interface{}),
		lock: sync.Mutex{},
		game: game{
			Players:   make([]Player, 10),
			turnStage: Uninitialized,
			// deck:        newDeck(),   // TODO init later on
			// lastElected: Utils.Set{}, // TODO init later on
			// executed:    Utils.Set{}, // TODO init later on
		},
	}
	go G.Handler()
	return &G
}

func (G *Game) Handler() {
	g := G.game
	vetoVoted := Utils.Set{}
	vetoResult := uint8(0)
	defer close(G.out)
	defer close(G.in)
	for {
		in := <-G.in
		switch in.(type) {
		// Join the game
		case Join:
			if G.game.turnStage == Uninitialized {
				id := in.(Join)
				p := Player{
					id: *id,
				}
				if len(g.Players) > 10 {
					G.out <- Error(LobbyFull)
					continue
				}
				g.PlayersMap[*id] = &p
				g.Players = append(g.Players, p)
				G.out <- Ack(General)
			} else {
				G.out <- Error(WrongPhase)
			}
		// Start the game
		case Start:
			if g.turnStage == Uninitialized {
				switch len(g.Players) {
				case 5, 6:
					g.setRoles(1)
				case 7, 8:
					g.setRoles(2)
				case 9, 10:
					g.setRoles(3)
				default:
					G.out <- Error(NotEnoughPlayers)
					continue
				}
				g.newPresident()
				G.out <- Ack(General)
			} else {
				G.out <- Error(WrongPhase)
			}
		case ElectChancellor:
			e := in.(ElectChancellor)
			if g.turnStage == ChancellorNeeded && g.President.id == *e.caller{
				if v, ok := g.PlayersMap[*e.proposal]; ok && !g.lastElected.Has(*e.proposal) {
					g.Chancellor = v
					g.Votes = make(map[*Player]bool, 10)
					g.turnStage = GovernmentElection
					G.out <- Ack(General)
				} else {
					G.out <- Error(InvalidInput)
				}
			} else {
				G.out <- Error(WrongPhase)
			}
		case VoteGovernment:
			if g.turnStage == GovernmentElection {
				e := in.(VoteGovernment)
				f := G.governmentCastVote(*e.caller, e.vote)
				switch f {
				case VoteAck:
					// todo ack the vote in output
				case VoteError:
					// todo send error (invalid vote)
				case Reject:
					// todo ack the vote in output
					// todo warn that the vote has failed
					g.newPresident()
				case RejectAndForce:
					g.ElectionTracker = 0       	   // reset the tracker
					g.policyChoice = g.deck.draw(1) // draw the policy
					g.enactPolicyUnsafe() // todo handle return types and win conditions

					// todo ack the vote in output
					// todo warn that the vote has failed
					g.newPresident()
				case Pass:
					// todo ack the vote in output
					g.turnStage = PresidentPolicies
					// todo warn that the vote has passed
				}
			} else {
				G.out <- Error(WrongPhase)
			}
		case PolicyChoice:
			e := in.(PolicyChoice)
			switch G.policyChoice(e.caller, e.selected) {
			case PolicyAck:
				switch g.turnStage {
				case ChancellorPolicies:
					// todo send confirmation of vote and chancellor vote request
				case VetoVote:
					// todo send confirmation of vote and veto request
				case VoteEnaction:
					// todo send confirmation of vote
					G.enactPolicy() // todo handle return types and win conditions
				}
			case PolicyError:
				// todo send error
			}
		case PolicyVeto:
			if g.turnStage == VetoVote {
				e := in.(PolicyVeto)
				if !vetoVoted.Has(*e.caller) && (g.PlayersMap[*e.caller] == g.President || g.PlayersMap[*e.caller] == g.
					Chancellor) {
					vetoVoted.Add(*e.caller)
					vetoResult++
					// todo add veto ack
					if len(vetoVoted) == 2 {
						if vetoResult == 2 {
							G.lock.Lock()
							// election tracker forces policies
							if g.ElectionTracker == 2 {
								g.ElectionTracker = 0       	   // reset the tracker
								g.policyChoice = g.deck.draw(1) // draw the policy
								g.enactPolicyUnsafe() // todo handle win conditions
							}
							g.ElectionTracker++
							g.newPresident()
							G.lock.Unlock()
							// todo warn that the government vetoed
						} else {
							G.enactPolicy()
							g.turnStage = SpecialPower
							// todo warn that the government failed the veto vote
						}
					} else {
						//todo send veto received
					}
				} else {
					// todo send error
				}
			}
		}
	}
}