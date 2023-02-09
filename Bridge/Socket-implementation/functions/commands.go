package functions

import (
	"encoding/json"
	"fmt"

	"github.com/slack-go/slack"
)

var ACTION_IDS = [...]string{"submitleave", "submit_options_insights"}
var ACTIONS = [...]string{"/help", "/apply-leave", "/cancel-leave", "/insights"}

func GetActionIds() []string {
	return ACTION_IDS[:]
}

// func createSectionBlock(text string) slack.Block {
// 	return slack.NewSectionBlock(
// 		&slack.TextBlockObject{
// 			Type: slack.MarkdownType,
// 			Text: text,
// 		},
// 		nil,
// 		nil,
// 	)
// }
// func createActionBlock(text string, action_id string) slack.Block {
// 	return slack.NewActionBlock(
// 		action_id,
// 		createButton(text, action_id),
// 	)
// }
// func createButton(text string, action_id string) *slack.ButtonBlockElement {
// 	return slack.NewButtonBlockElement(
// 		"",
// 		action_id,
// 		&slack.TextBlockObject{
// 			Type: slack.PlainTextType,
// 			Text: text,
// 		},
// 	)
// }
// func createInputBlock(text string, action_id string) slack.Block {
// 	return slack.NewInputBlock(
// 		action_id,
// 		&slack.TextBlockObject{
// 			Type: slack.PlainTextType,
// 			Text: text,
// 		},
// 		createDatePicker(text, action_id),
// 	)
// }

// func createDatePicker(text string, action_id string) map[string]interface{} {

// 	placeholder := map[string]interface{}{
// 		"type":  "plain_text",
// 		"text":  text,
// 		"emoji": true,
// 	}
// 	return map[string]interface{}{
// 		"type":         "datepicker",
// 		"initial_date": "1990-04-28",
// 		"placeholder":  placeholder,
// 		"action_id":    action_id,
// 	}
// }
// func createBlockElements(elements []slack.BlockElement) []slack.BlockElement {
// 	return slack.BlockElements{
// 		Elements: elements,
// 	}
// }

func testCommand() (response map[string]interface{}) {
	response = map[string]interface{}{
		"blocks": []slack.Block{
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: "foo",
				},
				nil,
				slack.NewAccessory(
					slack.NewButtonBlockElement(
						"",
						"somevalue",
						&slack.TextBlockObject{
							Type: slack.PlainTextType,
							Text: "bar",
						},
					),
				),
			),
		},
	}
	return response
}
func getRangePicker() (fromDate *slack.DatePickerBlockElement, toDate *slack.DatePickerBlockElement) {
	fromDate = slack.NewDatePickerBlockElement("from_date")
	fromDate.Placeholder = slack.NewTextBlockObject("plain_text", "From Date", false, false)
	toDate = slack.NewDatePickerBlockElement("to_date")
	toDate.Placeholder = slack.NewTextBlockObject("plain_text", "To Date", false, false)
	return fromDate, toDate
}
func leaveCommand() map[string]interface{} {
	// TODO:create a date picker function
	fromDate, toDate := getRangePicker()
	submitButton := slack.NewButtonBlockElement("submitleave", "submit", slack.NewTextBlockObject("plain_text", "Submit", false, false))
	response := map[string]interface{}{
		"blocks": []slack.Block{
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: "Apply Leave",
				},
				nil,
				nil,
			),
			slack.NewActionBlock("leave", fromDate, toDate, submitButton),
		},
	}
	fmt.Println(response)
	return response
}
func helpCommand() map[string]interface{} {
	blocks := []slack.Block{}
	blocks = append(blocks, slack.NewSectionBlock(
		&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: "Help",
		},
		nil,
		nil,
	),
	)
	blocks = append(blocks, slack.NewDividerBlock())
	for i := range ACTIONS {
		blocks = append(blocks, slack.NewSectionBlock(
			&slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: ACTIONS[i],
			},
			nil,
			nil,
		))
	}
	response := map[string]interface{}{
		"blocks": blocks,
	}
	return response
}
func cancelLeaveCommand(cmd *slack.SlashCommand) map[string]interface{} {
	user_id := cmd.UserID
	SendCancelLeave(user_id)
	return map[string]interface{}{
		"blocks": []slack.Block{
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: "applied for leave cancellation",
				},
				nil,
				nil,
			),
		},
	}
}

func insightsCommand() map[string]interface{} {
	fromDate, toDate := getRangePicker()
	options := []slack.OptionBlockObject{
		{
			Text:  slack.NewTextBlockObject("plain_text", "PDF", false, false),
			Value: "pdf",
		},
		{
			Text:  slack.NewTextBlockObject("plain_text", "CSV", false, false),
			Value: "csv",
		},
	}
	staticSelect := slack.NewOptionsSelectBlockElement("static_select", slack.NewTextBlockObject("plain_text", "Select Format Type", false, false), "format", &options[0], &options[1])

	submitButton := slack.NewButtonBlockElement("submit_options_insights", "submit", slack.NewTextBlockObject("plain_text", "Submit", false, false))
	return map[string]interface{}{
		"blocks": []slack.Block{
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: "Insights",
				},
				nil,
				nil,
			),
			slack.NewActionBlock("insights", fromDate, toDate, staticSelect, submitButton),
		},
	}
}

func ExecuteCommand(cmd *slack.SlashCommand) map[string]interface{} {
	command := cmd.Command
	var response map[string]interface{}
	switch command {
	case "/test":
		response = testCommand()
	case "/apply-leave":
		response = leaveCommand()
	case "/cancel-leave":
		response = cancelLeaveCommand(cmd)
	case "/insights":
		response = insightsCommand()
	case "/help":
		response = helpCommand()
	default:
		fmt.Println("Unknown command")
	}
	return response
}

func ProcessAction(callback *slack.InteractionCallback) slack.Attachment {
	var response slack.Attachment
	switch callback.ActionCallback.BlockActions[0].BlockID {
	case "leave":
		response = processLeave(callback)
	case "insights":
		response = processInsights(callback)
	default:
		fmt.Println("Unknown action")
	}
	return response
}

func processLeave(callback *slack.InteractionCallback) slack.Attachment {
	var stateofleave LeaveState
	user_id := callback.User.ID
	action := callback.ActionCallback.BlockActions[0].ActionID
	json.Unmarshal(callback.RawState, &stateofleave)
	fmt.Println(stateofleave)
	if action == "submitleave" {
		SendLeave(user_id, stateofleave.Values.Leave.FromDate.SelectedDate, stateofleave.Values.Leave.ToDate.SelectedDate)
	}
	return slack.Attachment{
		Text: "Leave applied",
	}
}

func processInsights(callback *slack.InteractionCallback) slack.Attachment {
	var stateofinsights InsightsState
	responseText := "You are not authorized to view this"
	user_id := callback.User.ID
	action := callback.ActionCallback.BlockActions[0].ActionID
	json.Unmarshal(callback.RawState, &stateofinsights)
	fmt.Println(stateofinsights)
	if action == "submit_options_insights" {
		responseText = GetInsights(user_id, stateofinsights.Values.Insights.FromDate.SelectedDate, stateofinsights.Values.Insights.ToDate.SelectedDate, stateofinsights.Values.Insights.Format.Option.Value)
	}
	return slack.Attachment{
		Text: responseText,
	}
}
