package vote

type HandlerFuncEnum string

const (
	HandlerFuncCreateVote                  HandlerFuncEnum = "CreateVote"
	HandlerFuncGetTotalVotes               HandlerFuncEnum = "GetTotalVotes"
	HandlerFuncGetTotalVotesForParticipant HandlerFuncEnum = "GetTotalVotesForParticipant"
	HandlerFuncGetTotalVotesForHour        HandlerFuncEnum = "GetTotalVotesForHour"
	HandlerFuncGetWinner                   HandlerFuncEnum = "GetWinner"
	HandlerFuncGetVotesFromParticipant     HandlerFuncEnum = "GetVotesFromParticipant"
)

func (h HandlerFuncEnum) String() string {
	return string(h)
}
