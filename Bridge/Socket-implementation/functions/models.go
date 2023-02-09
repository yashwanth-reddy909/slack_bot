package functions

// "state": {
//       "values": {
//         "leave": {
//           "from_date": {
//             "type": "datepicker",
//             "selected_date": "2021-01-02"
//           },
//           "to_date": {
//             "type": "datepicker",
//             "selected_date": "2021-01-01"
//           }
//         }
//       }
//     }
type datepickerstate struct {
	Type         string `json:"type"`
	SelectedDate string `json:"selected_date"`
}
type leavefieldState struct {
	FromDate datepickerstate `json:"from_date"`
	ToDate   datepickerstate `json:"to_date"`
}

type LeaveState struct {
	Values struct {
		Leave leavefieldState `json:"leave"`
	} `json:"values"`
}

// "state": {
//       "values": {
//         "insights": {
//           "from_date": {
//             "type": "datepicker",
//             "selected_date": "2023-02-09"
//           },
//           "to_date": {
//             "type": "datepicker",
//             "selected_date": "2023-02-09"
//           },
//           "format": {
//             "type": "static_select",
//             "selected_option": {
//               "text": {
//                 "type": "plain_text",
//                 "text": "PDF",
//                 "emoji": true
//               },
//               "value": "pdf"
//             }
//           }
//         }
//       }
//     }
type staticselectstate struct {
	Type  string `json:"type"`
	Option struct {
		Text struct {
			Type  string `json:"type"`
			Text  string `json:"text"`
			Emoji bool   `json:"emoji"`
		} `json:"text"`
		Value string `json:"value"`
	} `json:"selected_option"`
}
type insightsfieldState struct {
	FromDate datepickerstate `json:"from_date"`
	ToDate   datepickerstate `json:"to_date"`
	Format   staticselectstate `json:"format"`
}

type InsightsState struct {
	Values struct {
		Insights insightsfieldState `json:"insights"`
	} `json:"values"`
}