package errors

import "errors"

var (
	ErrorTeamExists  = errors.New("TEAM_EXISTS")
	ErrorPRExists    = errors.New("PR_EXISTS")
	ErrorPRMerged    = errors.New("PR_MERGED")
	ErrorNotAssigned = errors.New("NOT_ASSIGNED")
	ErrorNoCandidate = errors.New("NO_CANDIDATE")
	ErrorNotFound    = errors.New("NOT_FOUND")
	ErrInternal      = errors.New("INTERNAL")
)
