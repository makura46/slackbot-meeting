package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"reflect"
)

var BotRoom = "D7YU9FK0R"	// BotのChannel

type Room struct {
	UserName	map[string]bool
	Count	int
	People	int
	Status	int
}

type Controll map[string]Room

var botToken = ""	// Your App Token

func main() {
	api := slack.New(botToken)
	con := map[string]*Room{}
	person := map[string]string{}
	con[BotRoom] = &Room{map[string]bool{}, 0, 0, 1}

	api.SetDebug(true)
	initMsg := "参加する人は発言をしてください\n全員発言したら「start」で会議を始めます"
	startMsg := "会議を開始します\n終了するときは「finish」と発言してください"

	finishMsg := "会議を終了します"
	stopMsg := "会議が開始されてません"

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Println("Event Received: ")
		fmt.Println(reflect.TypeOf(msg.Data))
		switch  ev := msg.Data.(type) {
		case *slack.HelloEvent:
		case *slack.AckMessage:
			fmt.Println(ev)
		case *slack.ConnectedEvent:

		case *slack.MessageEvent:
			roomName := ev.Channel
			fmt.Println(roomName)
			if _, ok := con[ev.Channel]; ok {
				status := con[ev.Channel].Status
				if ev.Text == "start" && status == 0{
					con[ev.Channel].Status = 1
					rtm.SendMessage(rtm.NewOutgoingMessage(startMsg, ev.Channel))
				} else if ev.Text == "finish" {
					if status == 0 {
						rtm.SendMessage(rtm.NewOutgoingMessage(stopMsg, ev.Channel))
					} else if status == 1 {
						rtm.SendMessage(rtm.NewOutgoingMessage(finishMsg, ev.Channel))
						con[ev.Channel].Status = 2
					}
				}
				status = con[ev.Channel].Status
				switch status {
				case 0:
					if _, o := con[ev.Channel].UserName[ev.User]; !o {
						con[ev.Channel].People++
						con[ev.Channel].UserName[ev.User] = false
						person[ev.User] = ev.Channel
					}
				case 1:
					fmt.Println()
					fmt.Println(roomName)
					fmt.Println(BotRoom)
					fmt.Println()
					if roomName != BotRoom {
						break
					}
					var channel string
					var ok bool
					var p, c float64
					if channel, ok = person[ev.User]; ok && !con[channel].UserName[ev.User] {
						con[channel].UserName[ev.User] = true
						con[channel].Count++
						p = float64(con[channel].People)
						c = float64(con[channel].Count)
					}
					if p != 0 && c/p*100 >= 50 {
						fmt.Println("50%以上だよ")
						for key, _:= range con[channel].UserName {
							con[channel].UserName[key] = false
						}
						con[channel].Count = 0
					}
				}
			} else {
				rtm.SendMessage(rtm.NewOutgoingMessage(initMsg, ev.Channel))
				con[ev.Channel] = &Room{UserName: map[string]bool{}, Count: 0, People: 0,  Status: 0}
			}


		case *slack.PresenceChangeEvent:

		case *slack.LatencyReport:

		case *slack.RTMError:

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

		}
	}
}
