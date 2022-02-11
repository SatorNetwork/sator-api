package two_players

import "github.com/SatorNetwork/sator-api/svc/quiz_v2/message"

var (
	defaultUserExpectedMessages = []*message.Message{
		{
			MessageType:           message.PlayerIsJoinedMessageType,
			PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{},
		},
		{
			MessageType:           message.PlayerIsJoinedMessageType,
			PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{},
		},
		{
			MessageType: message.CountdownMessageType,
			CountdownMessage: &message.CountdownMessage{
				SecondsLeft: 3,
			},
		},
		{
			MessageType: message.CountdownMessageType,
			CountdownMessage: &message.CountdownMessage{
				SecondsLeft: 2,
			},
		},
		{
			MessageType: message.CountdownMessageType,
			CountdownMessage: &message.CountdownMessage{
				SecondsLeft: 1,
			},
		},
		{
			MessageType: message.CountdownMessageType,
			CountdownMessage: &message.CountdownMessage{
				SecondsLeft: 0,
			},
		},
		{
			MessageType: message.QuestionMessageType,
			QuestionMessage: &message.QuestionMessage{
				QuestionText:   "Joey played Dr. Drake Ramoray on which soap opera show?",
				TimeForAnswer:  1,
				QuestionNumber: 0,
				AnswerOptions: []message.AnswerOption{
					{
						AnswerText: "Santa Barbara",
					},
					{
						AnswerText: "Neighbours",
					},
					{
						AnswerText: "General Hospital",
					},
					{
						AnswerText: "Days of Our Lives",
					},
				},
			},
		},
		{
			MessageType: message.AnswerReplyMessageType,
			AnswerReplyMessage: &message.AnswerReplyMessage{
				Success:    true,
				SegmentNum: 1,
			},
		},
		{
			MessageType: message.QuestionMessageType,
			QuestionMessage: &message.QuestionMessage{
				QuestionText:   "What store does Phoebe hate?",
				TimeForAnswer:  1,
				QuestionNumber: 1,
				AnswerOptions: []message.AnswerOption{
					{
						AnswerText: "Amazon",
					},
					{
						AnswerText: "Costco",
					},
					{
						AnswerText: "Pottery Barn",
					},
					{
						AnswerText: "Walmart",
					},
				},
			},
		},
		{
			MessageType: message.AnswerReplyMessageType,
			AnswerReplyMessage: &message.AnswerReplyMessage{
				Success:    true,
				SegmentNum: 1,
			},
		},
		{
			MessageType: message.QuestionMessageType,
			QuestionMessage: &message.QuestionMessage{
				QuestionText:   "Rachel got a job with which company in Paris?",
				TimeForAnswer:  1,
				QuestionNumber: 2,
				AnswerOptions: []message.AnswerOption{
					{
						AnswerText: "Cartier",
					},
					{
						AnswerText: "Gucci",
					},
					{
						AnswerText: "Louis Vuitton",
					},
					{
						AnswerText: "Zara",
					},
				},
			},
		},
		{
			MessageType: message.AnswerReplyMessageType,
			AnswerReplyMessage: &message.AnswerReplyMessage{
				Success:    true,
				SegmentNum: 1,
			},
		},
		{
			MessageType: message.QuestionMessageType,
			QuestionMessage: &message.QuestionMessage{
				QuestionText:   "Phoebe’s scientist boyfriend David worked in what city?",
				TimeForAnswer:  1,
				QuestionNumber: 3,
				AnswerOptions: []message.AnswerOption{
					{
						AnswerText: "Minsk",
					},
					{
						AnswerText: "Moscow",
					},
					{
						AnswerText: "Kyiv",
					},
					{
						AnswerText: "Berlin",
					},
				},
			},
		},
		{
			MessageType: message.AnswerReplyMessageType,
			AnswerReplyMessage: &message.AnswerReplyMessage{
				Success:    true,
				SegmentNum: 1,
			},
		},
		{
			MessageType: message.QuestionMessageType,
			QuestionMessage: &message.QuestionMessage{
				QuestionText:   "Monica dated an ophthalmologist named?",
				TimeForAnswer:  1,
				QuestionNumber: 4,
				AnswerOptions: []message.AnswerOption{
					{
						AnswerText: "Chandler",
					},
					{
						AnswerText: "Rafael",
					},
					{
						AnswerText: "Richard",
					},
					{
						AnswerText: "Robert",
					},
				},
			},
		},
		{
			MessageType: message.AnswerReplyMessageType,
			AnswerReplyMessage: &message.AnswerReplyMessage{
				Success:    true,
				SegmentNum: 1,
			},
		},
		{
			MessageType:         message.WinnersTableMessageType,
			WinnersTableMessage: &message.WinnersTableMessage{},
		},
	}
)
