package entity

type Round struct {
	ID           string
	Nome         string
	Participants []Participant
	CreatedAt    int64
}

type Participant struct {
	ID   string
	Nome string
}

type Vote struct {
	RoundID       string
	ParticipantID string
	Timestamp     int64
	IP            string
}
