package github.com/evoteum/planzoco/go/planzoco/models

type Event struct {
	ID        string     `json:"id"`
	Name      string     `json:"name" form:"name" binding:"required"`
	Questions []Question `json:"questions,omitempty"`
}

type Question struct {
	ID      string   `json:"id"`
	EventID string   `json:"event_id"`
	Text    string   `json:"text" form:"text" binding:"required"`
	Options []Option `json:"options,omitempty"`
}

func (q Question) WinningOptions() []Option {
	if len(q.Options) == 0 {
		return nil
	}

	allZero := true
	maxVotes := q.Options[0].Votes

	for _, opt := range q.Options {
		if opt.Votes > 0 {
			allZero = false
		}
		if opt.Votes > maxVotes {
			maxVotes = opt.Votes
		}
	}

	if allZero {
		return nil
	}

	// Collect all options with max votes
	var winners []Option
	for _, opt := range q.Options {
		if opt.Votes == maxVotes {
			winners = append(winners, opt)
		}
	}

	return winners
}

type Option struct {
	ID         string `json:"id"`
	QuestionID string `json:"question_id"`
	Text       string `json:"text" form:"text" binding:"required"`
	Votes      int    `json:"votes"`
}
