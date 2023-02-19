package dblayer

import "time"

type Iuser struct {
	userID       int32     //`bson:"userID"`
	userFname    string    //`bson:"userFname"`
	userLname    string    //`bson:"userLname"`
	userMobile   string    //`bson:"userMobile"`
	userAddress  string    //`bson:"userAddress"`
	userCodeMeli string    //`bson:"userCodeMeli"`
	userDateTime time.Time //`bson:"userDateTime"`
	userIsActive bool      //`bson:"userIsActive"`
}

func NewUser(userID int32, userFname string, userLname string, userMobile string, userAddress string, userCodeMeli string, userDateTime time.Time, userIsActive bool) *Iuser {
	return &Iuser{
		userID:       userID,
		userFname:    userFname,
		userLname:    userLname,
		userMobile:   userMobile,
		userAddress:  userAddress,
		userCodeMeli: userCodeMeli,
		userDateTime: userDateTime,
		userIsActive: userIsActive,
	}
}
