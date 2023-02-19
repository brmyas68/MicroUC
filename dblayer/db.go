package dblayer

import (
	"log"

	"context"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUCDataBaseService interface {
	ExistsUser(userMobile string) int32
	InsertToken(token *IToken) bool
	DeleteToken(HashToken string) bool
	GetRoleID(userID int32) int32
	GetRoleName(roleID int32) string
	GetUserbyToken(HashToken string) map[string]interface{} // return Iuser
	GetAllUsers() []bson.M                                  // return Iuser
	GenerateTokenLastID() int32
	InsertUser(user *Iuser) bool
	GenerateUserLastID() int32
	InsertUserRole(userRole *IUserRole) bool
	GenerateUserRoleLastID() int32
}

type UCDataBaseServiceStruct struct {
	db *mongo.Client
}

func NewUCDataBaseServiceStruct(db *mongo.Client) IUCDataBaseService {
	return &UCDataBaseServiceStruct{db: db}
}

func (tsUC *UCDataBaseServiceStruct) GenerateTokenLastID() int32 {
	dbToken := tsUC.db.Database("db_uc")
	dbCollection := dbToken.Collection("tbl_token")

	total, err := dbCollection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		//log.Fatal(" not found : ", err)
		return 1
	}
	return int32(total + 1)
}

func (tsUC *UCDataBaseServiceStruct) GenerateUserLastID() int32 {
	dbToken := tsUC.db.Database("db_uc")
	dbCollection := dbToken.Collection("tbl_user")

	total, err := dbCollection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		//log.Fatal(" not found : ", err)
		return 1
	}
	return int32(total + 1)
}

func (tsUC *UCDataBaseServiceStruct) GenerateUserRoleLastID() int32 {
	dbToken := tsUC.db.Database("db_uc")
	dbCollection := dbToken.Collection("tbl_userrole")

	total, err := dbCollection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		//log.Fatal(" not found : ", err)
		return 1
	}
	return int32(total + 1)
}

func (tsUC *UCDataBaseServiceStruct) ExistsUser(userMobile string) int32 {
	dbUser := tsUC.db.Database("db_uc")
	dbCollection := dbUser.Collection("tbl_user")

	var entity bson.M
	err := dbCollection.FindOne(context.Background(), bson.D{{Key: "userMobile", Value: userMobile}}).Decode(&entity)
	if err != nil {
		return 0
	}

	UsrID, err := strconv.ParseInt(fmt.Sprint(entity["userID"]), 10, 64)
	if err != nil {
		panic(err)
	}
	return int32(UsrID)
}

func (tsUC *UCDataBaseServiceStruct) InsertToken(token *IToken) bool {
	dbToken := tsUC.db.Database("db_uc")
	dbCollection := dbToken.Collection("tbl_token")
	Result, err := dbCollection.InsertOne(context.Background(), bson.D{
		{Key: "tokenID", Value: token.tokenID},
		{Key: "tokenUserID", Value: token.tokenUserID},
		{Key: "tokenHash", Value: token.tokenHash},
		{Key: "tokenDateTime", Value: token.tokenDateTime},
	})
	if err != nil {
		fmt.Println(" Not Insert Record")
	}
	if Result.InsertedID != nil {
		fmt.Printf("%v, type = %T\n", Result.InsertedID, Result.InsertedID)
		return true
	}
	return false
}

func (tsUC *UCDataBaseServiceStruct) GetRoleID(userID int32) int32 {
	dbUserRole := tsUC.db.Database("db_uc")
	dbCollection := dbUserRole.Collection("tbl_userRole")

	var entity bson.M
	err := dbCollection.FindOne(context.Background(), bson.D{{Key: "userRoleUserID", Value: userID}}).Decode(&entity)
	if err != nil {
		return 0
	}

	userRoleRoleID, err := strconv.ParseInt(fmt.Sprint(entity["userRoleRoleID"]), 10, 64)
	if err != nil {
		return 0
	}
	return int32(userRoleRoleID)
}

func (tsUC *UCDataBaseServiceStruct) GetRoleName(roleID int32) string {
	dbRole := tsUC.db.Database("db_uc")
	dbCollection := dbRole.Collection("tbl_Role")

	var entity bson.M
	err := dbCollection.FindOne(context.Background(), bson.D{{Key: "roleID", Value: roleID}}).Decode(&entity)
	if err != nil {
		log.Fatal(" Not  Found", err)
	}

	roleName := fmt.Sprint(entity["roleName"])
	return roleName
}

func (tsUC *UCDataBaseServiceStruct) GetUserbyToken(HashToken string) map[string]interface{} {
	dbUc := tsUC.db.Database("db_uc")
	dbCollectionToken := dbUc.Collection("tbl_token")
	dbCollectionUser := dbUc.Collection("tbl_user")

	var entityToken bson.M
	err := dbCollectionToken.FindOne(context.Background(), bson.D{{Key: "tokenHash", Value: HashToken}}).Decode(&entityToken)
	if err != nil {
		log.Fatal(" not found : ", err)
	}

	tokenUserID, err := strconv.ParseInt(fmt.Sprint(entityToken["tokenUserID"]), 10, 64)
	if err != nil {
		log.Fatal(" not found : ", err)
	}

	var entityUser bson.M

	errusr := dbCollectionUser.FindOne(context.Background(), bson.D{{Key: "userID", Value: tokenUserID}}).Decode(&entityUser)
	if errusr != nil {
		log.Fatal(" not found : ", errusr)
	}

	return entityUser
}

func (tsUC *UCDataBaseServiceStruct) GetAllUsers() []bson.M {
	dbUc := tsUC.db.Database("db_uc")
	dbCollectionUser := dbUc.Collection("tbl_user")

	Result, err := dbCollectionUser.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(" not found : ", err)
	}

	defer Result.Close(context.Background())

	var Results []bson.M
	if err = Result.All(context.TODO(), &Results); err != nil {
		panic(err)
	}

	return Results
}

func (tsUC *UCDataBaseServiceStruct) DeleteToken(HashToken string) bool {
	dbUc := tsUC.db.Database("db_uc")
	dbCollectionToken := dbUc.Collection("tbl_token")

	Result, err := dbCollectionToken.DeleteOne(context.Background(), bson.D{{Key: "tokenHash", Value: HashToken}})
	if err != nil {
		log.Fatal(" not found : ", err)
	}

	if Result.DeletedCount > 0 {
		return true
	}
	return false
}

func (tsUC *UCDataBaseServiceStruct) InsertUser(user *Iuser) bool {
	dbUc := tsUC.db.Database("db_uc")
	dbCollectionUser := dbUc.Collection("tbl_user")

	Result, err := dbCollectionUser.InsertOne(context.Background(), bson.D{
		{Key: "userID", Value: user.userID},
		{Key: "userFname", Value: user.userFname},
		{Key: "userLname", Value: user.userLname},
		{Key: "userMobile", Value: user.userMobile},
		{Key: "userCodeMeli", Value: user.userCodeMeli},
		{Key: "userIsActive", Value: user.userIsActive},
		{Key: "userAddress", Value: user.userAddress},
		{Key: "userDateTime", Value: user.userDateTime},
	})
	if err != nil {
		fmt.Println(" Not Insert Record")
	}
	if Result.InsertedID != nil {
		fmt.Printf("%v, type = %T\n", Result.InsertedID, Result.InsertedID)
		return true
	}
	return false

}

func (tsUC *UCDataBaseServiceStruct) InsertUserRole(userRole *IUserRole) bool {
	dbUc := tsUC.db.Database("db_uc")
	dbCollectionUserRole := dbUc.Collection("tbl_userrole")

	Result, err := dbCollectionUserRole.InsertOne(context.Background(), bson.D{
		{Key: "userRoleID", Value: userRole.userRoleID},
		{Key: "userRoleUserID", Value: userRole.userRoleUserID},
		{Key: "userRoleRoleID", Value: userRole.userRoleRoleID},
	})
	if err != nil {
		fmt.Println(" Not Insert Record")
	}
	if Result.InsertedID != nil {
		fmt.Printf("%v, type = %T\n", Result.InsertedID, Result.InsertedID)
		return true
	}
	return false

}
