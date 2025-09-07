package pipe

type ExecutionType string

const (
	SEQUENTIAL                     ExecutionType = "SEQUENTIAL"
	CONCURRENT                     ExecutionType = "CONCURRENT"
	SEQUENTIAL_WITH_FIRST_RESULT   ExecutionType = "SEQUENTIAL_WITH_FIRST_RESULT"
	SEQUENTIAL_BLOCKING_ONLY_FIRST ExecutionType = "SEQUENTIAL_BLOCKING_ONLY_FIRST"
)

func (e ExecutionType) String() string {
	return string(e)
}

func ParseExecutionType(s string) ExecutionType {
	switch s {
	case SEQUENTIAL.String():
		return SEQUENTIAL
	case CONCURRENT.String():
		return CONCURRENT
	case SEQUENTIAL_WITH_FIRST_RESULT.String():
		return SEQUENTIAL_WITH_FIRST_RESULT
	case SEQUENTIAL_BLOCKING_ONLY_FIRST.String():
		return SEQUENTIAL_BLOCKING_ONLY_FIRST
	default:
		return ""
	}
}
