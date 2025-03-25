package auth

import "time"

type TokenExpiredError struct {
	at time.Time
}

func (e TokenExpiredError) Error() string {
	return "token is expired at " + e.at.Format(time.RFC3339)
}

func NewTokenExpiredError(at time.Time) TokenExpiredError {
	return TokenExpiredError{at: at}
}
