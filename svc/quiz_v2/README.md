# Quiz service

## Messages

### Message types list

**Server side messages:**
- player_is_joined
- countdown
- question
- answer_reply
- winners_table

**Client side message**
- answer

### Message examples

**player_is_joined**
```json
{
  "message_type": 0,
  "player_is_joined_message": {
    "player_id": "080319a9-82d1-4b3f-b3a0-5c27f8a53b58",
    "username": "johndoe8011668518469047223"
  }
}
```

**countdown**
```json
{
  "message_type": 1,
  "countdown_message": {
    "SecondsLeft": 3
  }
}
```

**question**
```json
{
  "message_type": 2,
  "question_message": {
    "question_id": "d257ab4a-f47b-4fca-a019-9552374c4761",
    "question_text": "Joey played Dr. Drake Ramoray on which soap opera show?",
    "time_for_answer": 0,
    "question_number": 0,
    "answer_options": [
      {
        "answer_id": "dab76c4a-5b61-4610-a843-6f802203f351",
        "answer_text": "Days of Our Lives"
      },
      {
        "answer_id": "aa83cacb-2c81-439a-b2e1-81437f08d2bd",
        "answer_text": "General Hospital"
      },
      {
        "answer_id": "5730cde1-5a2e-484b-89b1-4184280239ad",
        "answer_text": "Santa Barbara"
      },
      {
        "answer_id": "54cf1b58-1bbf-4088-8f1c-7976eb6cc2b8",
        "answer_text": "Neighbours"
      }
    ]
  }
}

```

**answer**
```json
{
  "message_type": 3,
  "answer_message": {
    "user_id": "080319a9-82d1-4b3f-b3a0-5c27f8a53b58",
    "question_id": "d257ab4a-f47b-4fca-a019-9552374c4761",
    "answer_id": "dab76c4a-5b61-4610-a843-6f802203f351"
  }
}
```

**winners_table**
```json
{
  "message_type": 4,
  "answer_reply_message": {
    "success": true,
    "segment_num": 1,
    "is_fastest_answer": false
  }
}
```

**winners_table**
```json
{
  "message_type": 5,
  "winners_table_message": {
    "prize_pool_distribution": {
      "johndoe15872102387313413371": 120,
      "johndoe8011668518469047223": 130
    }
  }
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
