# Quiz service

## Messages

### Message types list

**Server side messages:**
- player_connected
- player_disconnected
- countdown
- question
- question_result
- challenge_result

**Client side message**
- answer


### Message structure

```json
{
    "type":"player_connected",
    "sent_at":"2021-05-28T18:42:57.969135+03:00",
    "payload": {
        "user_id":"24adf20e-61e1-439f-b650-56886320501e",
        "username":"johndoe",
    }
}
```

### Payload types

**player_connected / player_disconnected**
```json
{
    "user_id":"24adf20e-61e1-439f-b650-56886320501e",
    "username":"johndoe",
}
```

**countdown**
```json
{
    "countdown": 3 // from 3 to 0
}
```

**question**
```json
{
    "question_id": "24adf20e-61e1-439f-b650-56886320501e",
    "question_text": "Question text",
    "time_for_answer": 8,
    "total_questions": 10,
    "question_number": 3, // number of current question
    "answer_options": [
        {
            "answer_id": "14adf20e-61e1-439f-b650-56886320501e",
            "answer_text": "Answer text"
        },
        {
            "answer_id": "24adf20e-61e1-439f-b650-56886320501e",
            "answer_text": "Socond text"
        },
        {
            "answer_id": "34adf20e-61e1-439f-b650-56886320501e",
            "answer_text": "third text"
        },
        {
            "answer_id": "44adf20e-61e1-439f-b650-56886320501e",
            "answer_text": "fourth text"
        }
    ]
}
```

**answer**
```json
{
    "question_id": "24adf20e-61e1-439f-b650-56886320501e",
    "answer_id": "44adf20e-61e1-439f-b650-56886320501e"
}
```

**question_result**
```json
{
    "question_id": "24adf20e-61e1-439f-b650-56886320501e",
    "result": true,
    "rate": 3, // from 0 to 3
    "correct_answer_id": "44adf20e-61e1-439f-b650-56886320501e",
    "questions_left": 7,
    "additional_pts": 2,
}
```

**challenge_result**
```json
{
    "challenge_id": "44adf20e-61e1-439f-b650-56886320501e",
    "prize_pool": "250 SAO",
    "show_transaction_url": "https://....",
    "winners": [
        {
            "user_id": "24adf20e-61e1-439f-b650-56886320501e",
            "username": "johndoe",
            "prize": "2.45",
        }
    ]
}
```

---

## Reward example calculation:

| Prize pool | Total questions | Winners |
| ---------- | --------------- | ------- |
| 250 SAO    | 10              | 5       |
|            |                 |         |

### Formula:
**Rate**
```
rate = round(4-(round(answerTimeSec / (questionTimeSec / 4))))
```

**Reward**
```
reward = (prizePool / ((totalWinners * totalQuestions) + totalPts + totalRate)) * (totalQuestions + pts + rate)
```

| Winner # | PTS | Rate | Points | Reward |
| -------- | --- | ---- | ------ | ------ |
| #1       | 0   |      | 10     | ~21.74 |
| #2       | 10  |      | 20     | ~43.48 |
| #3       | 20  |      | 30     | ~65.22 |
| #4       | 30  |      | 40     | ~86.96 |
| #5       | 5   |      | 15     | ~32.61 |
|          |     |      |        |