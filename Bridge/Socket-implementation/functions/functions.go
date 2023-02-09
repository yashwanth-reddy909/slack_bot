package functions

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"golang.org/x/exp/slices"
)

func UserQA(client *socketmode.Client, questions []string, Answers map[string][]string, botId string) {

	for evt := range client.Events {
		// fmt.Println("************************************")

		switch evt.Type {
		case socketmode.EventTypeConnecting:
			fmt.Println("Connecting to Slack with Socket Mode...")
		case socketmode.EventTypeConnectionError:
			fmt.Println("Connection failed. Retrying later...")
		case socketmode.EventTypeConnected:
			fmt.Println("Connected to Slack with Socket Mode.")
		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				fmt.Printf("Ignored %+v\n", evt)
				continue
			}

			// fmt.Printf("Event received: %+v\n", eventsAPIEvent)

			client.Ack(*evt.Request)

			switch eventsAPIEvent.Type {
			case slackevents.CallbackEvent:
				innerEvent := eventsAPIEvent.InnerEvent
				switch ev := innerEvent.Data.(type) {
				case *slackevents.MessageEvent:
					if ev.User == botId {
						continue
					}

					if len(Answers[ev.User]) == len(questions) {
						client.Client.PostMessage(ev.User, slack.MsgOptionText("You have Answered all the questions", false))
						continue
					} else {
						fmt.Println(ev)
						Answers[ev.User] = append(Answers[ev.User], ev.Text)
						if len(Answers[ev.User]) == len(questions) {
							client.Client.PostMessage(ev.User, slack.MsgOptionText("You have Answered all the questions, Have a great day", false))
							fmt.Println("************************************")
							fmt.Println(ev)
							fmt.Println("************************************")
							SendResponseOfUserToScrum(client, questions, Answers, ev.User, "C04MC939RBR")

						} else {
							client.Client.SendMessage(ev.User, slack.MsgOptionText(questions[len(Answers[ev.User])], false))
						}
					}
				case *slackevents.AppMentionEvent:
					_, _, err := client.Client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
					if err != nil {
						fmt.Printf("failed posting message: %v", err)
					}
				case *slackevents.MemberJoinedChannelEvent:
					fmt.Printf("user %q joined to channel %q", ev.User, ev.Channel)
				}
			default:
				client.Debugf("unsupported Events API event received")
			}
		case socketmode.EventTypeInteractive:
			callback, ok := evt.Data.(slack.InteractionCallback)
			if !ok {
				fmt.Printf("Ignored %+v\n", evt)

				continue
			}

			// fmt.Printf("Interaction received: %+v\n", callback)
			client.Debugf("Interaction received")

			var payload interface{}

			switch callback.Type {
			case slack.InteractionTypeBlockActions:
				// See https://api.slack.com/apis/connections/socket-implement#button
				actionIds := GetActionIds()
				//check if the action id is in the list of action ids
				if slices.Contains(actionIds, callback.ActionCallback.BlockActions[0].ActionID) {
					attachment := ProcessAction(&callback)
					client.Debugf("button clicked!")
					client.Client.PostMessage(callback.Channel.ID, slack.MsgOptionAttachments(attachment))
				}

				client.Debugf("action received")
			case slack.InteractionTypeShortcut:
			case slack.InteractionTypeViewSubmission:
				// See https://api.slack.com/apis/connections/socket-implement#modal
			case slack.InteractionTypeDialogSubmission:
			default:

			}

			client.Ack(*evt.Request, payload)
		case socketmode.EventTypeSlashCommand:
			cmd, ok := evt.Data.(slack.SlashCommand)
			if !ok {
				fmt.Printf("Ignored %+v\n", evt)

				continue
			}

			client.Debugf("Slash command received: %+v", cmd.Command)

			// payload := map[string]interface{}{
			// 	"blocks": []slack.Block{
			// 		slack.NewSectionBlock(
			// 			&slack.TextBlockObject{
			// 				Type: slack.MarkdownType,
			// 				Text: "foo",
			// 			},
			// 			nil,
			// 			slack.NewAccessory(
			// 				slack.NewButtonBlockElement(
			// 					"",
			// 					"somevalue",
			// 					&slack.TextBlockObject{
			// 						Type: slack.PlainTextType,
			// 						Text: "bar",
			// 					},
			// 				),
			// 			),
			// 		),
			// 	}}
			response := ExecuteCommand(&cmd)
			client.Ack(*evt.Request, response)
		default:
			fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
		}

	}
	fmt.Println("End of QA")

}

func GetRandomColor() string {
	var colors []string
	colors = append(colors, "#ff6600")
	colors = append(colors, "#99ff33")
	colors = append(colors, "#ff6666")
	colors = append(colors, "#000066")
	colors = append(colors, "#ccff33")
	colors = append(colors, "#cc0066")
	colors = append(colors, "#00ccff")
	colors = append(colors, "#66ff99")
	colors = append(colors, "#3399ff")
	colors = append(colors, "#ff9933")
	index := rand.Intn(9)
	return colors[index]
}

func SendResponseOfUserToScrum(client *socketmode.Client, questions []string, Answers map[string][]string, userId string, channelId string) {
	username, image := getUserProfile(client, userId)
	var attachments []slack.Attachment

	for i := 0; i < len(questions); i++ {
		if Answers[userId][i] != "-" {
			var fields []slack.AttachmentField
			fields = append(fields, slack.AttachmentField{
				Title: questions[i],
				Value: Answers[userId][i],
			})
			attachments = append(attachments, slack.Attachment{

				Color:  GetRandomColor(),
				Fields: fields,
			})
		}
	}

	client.PostMessage(
		channelId,
		slack.MsgOptionText(username+" has posted an update for srum", false),
		slack.MsgOptionAttachments(attachments...),
		slack.MsgOptionUsername(username),
		slack.MsgOptionIconURL(image),
	)
}
func getUserProfile(client *socketmode.Client, userId string) (username string, image string) {
	user, err := client.Client.GetUserInfo(userId)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	return user.Profile.RealName, user.Profile.ImageOriginal
}
