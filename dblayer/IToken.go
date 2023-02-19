package dblayer

import "time"

type IToken struct {
	tokenID       int32
	tokenUserID   int32
	tokenHash     string
	tokenDateTime time.Time
}

func NewToken(tokenID int32, tokenUserID int32, tokenHash string, tokenDateTime time.Time) *IToken {
	return &IToken{
		tokenID:       tokenID,
		tokenUserID:   tokenUserID,
		tokenHash:     tokenHash,
		tokenDateTime: tokenDateTime,
	}
}
