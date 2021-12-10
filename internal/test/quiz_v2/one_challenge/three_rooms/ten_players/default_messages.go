package ten_players

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
			MessageType:           message.PlayerIsJoinedMessageType,
			PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{},
		},
		{
			MessageType:           message.PlayerIsJoinedMessageType,
			PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{},
		},
		{
			MessageType:           message.PlayerIsJoinedMessageType,
			PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{},
		},
		{
			MessageType:           message.PlayerIsJoinedMessageType,
			PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{},
		},
		{
			MessageType:           message.PlayerIsJoinedMessageType,
			PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{},
		},
		{
			MessageType:           message.PlayerIsJoinedMessageType,
			PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{},
		},
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
			MessageType: message.QuestionMessageType,
			QuestionMessage: &message.QuestionMessage{
				QuestionText:   "Joey played Dr. Drake Ramoray on which soap opera show?",
				TimeForAnswer:  0,
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
			MessageType:         message.WinnersTableMessageType,
			WinnersTableMessage: &message.WinnersTableMessage{},
		},
	}
)
