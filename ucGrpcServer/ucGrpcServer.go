package ucGrpcServer

import (
	"context"
	"fmt"
	"main/dblayer"
	"main/uc/pb"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UCGrpcServerStruct struct {
}

var UCDataBaseStruct dblayer.IUCDataBaseService

func NewUCGrpcServerStruct(db *mongo.Client) *UCGrpcServerStruct {

	UCDataBaseStruct = dblayer.NewUCDataBaseServiceStruct(db)
	return &UCGrpcServerStruct{}
}

func (tsUCGrpc *UCGrpcServerStruct) Login(ctx context.Context, Request *pb.RequestLogin) (*pb.ResponseLogin, error) {
	userID := UCDataBaseStruct.ExistsUser(Request.GetMobile())
	if userID > 0 {
		var UCTokenStruct dblayer.IGenerateTokenService = dblayer.NewGenerateTokenServiceStruct()
		token := UCTokenStruct.CreateToken("GrpcUC")
		hashToken := UCTokenStruct.CreateHash256(token)
		TokenID := UCDataBaseStruct.GenerateTokenLastID()
		timeToken := time.Now()
		IToken := dblayer.NewToken(TokenID, userID, hashToken, timeToken)
		stateToken := UCDataBaseStruct.InsertToken(IToken)
		if stateToken {
			return &pb.ResponseLogin{
				Token: token,
				Status: &pb.Status{
					StatusCode:    pb.StatusCode_Status200,
					StatusMessage: pb.StatusMessage_SUCCESS,
				},
			}, nil
		}
	}
	return &pb.ResponseLogin{
		Token: "",
		Status: &pb.Status{
			StatusCode:    pb.StatusCode_Status400,
			StatusMessage: pb.StatusMessage_FAILED,
		},
	}, nil
}

func (tsUCGrpc *UCGrpcServerStruct) Permission(ctx context.Context, Request *pb.RequestEmpty) (*pb.ResponsePermission, error) {

	if MD, OK := metadata.FromIncomingContext(ctx); OK {
		Token := MD.Get("authorization")[0]
		if Token != "" {
			var UCTokenStruct dblayer.IGenerateTokenService = dblayer.NewGenerateTokenServiceStruct()
			hashToken := UCTokenStruct.CreateHash256(Token)
			IUser := UCDataBaseStruct.GetUserbyToken(hashToken)
			UserID, err := strconv.ParseInt(fmt.Sprint(IUser["userID"]), 10, 64)
			if err != nil {
				panic(err)
			}
			RoleID := UCDataBaseStruct.GetRoleID(int32(UserID))
			if RoleID > 0 {
				RoleName := UCDataBaseStruct.GetRoleName(RoleID)
				if RoleName != "" && RoleName == "ADMIN" {
					return &pb.ResponsePermission{
						Status: &pb.Status{
							StatusCode:    pb.StatusCode_Status200,
							StatusMessage: pb.StatusMessage_ALLOW,
						}}, nil
				} else {
					return &pb.ResponsePermission{
						Status: &pb.Status{
							StatusCode:    pb.StatusCode_Status200,
							StatusMessage: pb.StatusMessage_DENY,
						}}, nil
				}
			}
		}

	}
	return &pb.ResponsePermission{
		Status: &pb.Status{
			StatusCode:    pb.StatusCode_Status400,
			StatusMessage: pb.StatusMessage_UNAUTHORIZED,
		},
	}, nil
}

func (tsUCGrpc *UCGrpcServerStruct) GetUser(ctx context.Context, Request *pb.RequestEmpty) (*pb.ResponseUser, error) {
	if MD, OK := metadata.FromIncomingContext(ctx); OK {
		Token := MD.Get("authorization")[0]
		if Token != "" {
			var UCTokenStruct dblayer.IGenerateTokenService = dblayer.NewGenerateTokenServiceStruct()
			hashToken := UCTokenStruct.CreateHash256(Token)
			IUser := UCDataBaseStruct.GetUserbyToken(hashToken)
			UserID, err := strconv.ParseInt(fmt.Sprint(IUser["userID"]), 10, 64)
			if err != nil {
				panic(err)
			}
			UsrFname := fmt.Sprint(IUser["userFname"])
			UsrLname := fmt.Sprint(IUser["userLname"])
			usrMobile := fmt.Sprint(IUser["userMobile"])
			usrAddress := fmt.Sprint(IUser["userAddress"])
			usrCodeMeli := fmt.Sprint(IUser["userCodeMeli"])
			usrDateTimeString, err := strconv.ParseInt(fmt.Sprint(IUser["userDateTime"]), 10, 64)
			if err != nil {
				panic(err)
			}
			usrDateTime := time.UnixMilli(usrDateTimeString).UTC()
			usrIsActive, err := strconv.ParseBool(fmt.Sprint(IUser["userIsActive"]))
			if err != nil {
				panic(err)
			}
			timestamp := timestamppb.New(usrDateTime)
			return &pb.ResponseUser{
				User: &pb.User{
					UserID:       int32(UserID),
					UserFname:    UsrFname,
					UserLname:    UsrLname,
					UserMobile:   usrMobile,
					UserAddress:  usrAddress,
					UserCodeMeli: usrCodeMeli,
					UserDateTime: timestamp,
					UserIsActive: usrIsActive,
				},
			}, nil

		}
	}
	return nil, nil
}

func (tsUCGrpc *UCGrpcServerStruct) GetAllUser(Request *pb.RequestEmpty, stream pb.UCService_GetAllUserServer) error {

	var Users = UCDataBaseStruct.GetAllUsers()
	for _, User := range Users {

		UserID, err := strconv.ParseInt(fmt.Sprint(User["userID"]), 10, 64)
		if err != nil {
			panic(err)
		}
		UsrFname := fmt.Sprint(User["userFname"])
		UsrLname := fmt.Sprint(User["userLname"])
		usrMobile := fmt.Sprint(User["userMobile"])
		usrAddress := fmt.Sprint(User["userAddress"])
		usrCodeMeli := fmt.Sprint(User["userCodeMeli"])
		usrDateTimeString, err := strconv.ParseInt(fmt.Sprint(User["userDateTime"]), 10, 64)
		if err != nil {
			panic(err)
		}
		usrDateTime := time.UnixMilli(usrDateTimeString).UTC()
		usrIsActive, err := strconv.ParseBool(fmt.Sprint(User["userIsActive"]))
		if err != nil {
			panic(err)
		}
		timestamp := timestamppb.New(usrDateTime)

		stream.Send(&pb.ResponseAllUser{
			User: &pb.User{
				UserID:       int32(UserID),
				UserFname:    UsrFname,
				UserLname:    UsrLname,
				UserMobile:   usrMobile,
				UserAddress:  usrAddress,
				UserCodeMeli: usrCodeMeli,
				UserDateTime: timestamp,
				UserIsActive: usrIsActive,
			},
			Status: &pb.Status{
				StatusCode:    pb.StatusCode_Status200,
				StatusMessage: pb.StatusMessage_SUCCESS,
			},
		})
	}

	return nil
}

func (tsUCGrpc *UCGrpcServerStruct) LogOut(ctx context.Context, Request *pb.RequestEmpty) (*pb.ResponseLogOut, error) {
	if MD, OK := metadata.FromIncomingContext(ctx); OK {
		Token := MD.Get("authorization")[0]
		if Token != "" {
			var UCTokenStruct dblayer.IGenerateTokenService = dblayer.NewGenerateTokenServiceStruct()
			hashToken := UCTokenStruct.CreateHash256(Token)
			State := UCDataBaseStruct.DeleteToken(hashToken)
			if State {
				return &pb.ResponseLogOut{
					Status: &pb.Status{
						StatusCode:    pb.StatusCode_Status200,
						StatusMessage: pb.StatusMessage_SUCCESS,
					}}, nil
			} else {
				return &pb.ResponseLogOut{
					Status: &pb.Status{
						StatusCode:    pb.StatusCode_Status400,
						StatusMessage: pb.StatusMessage_FAILED,
					}}, nil
			}
		}
	}
	return nil, nil
}

func (tsUCGrpc *UCGrpcServerStruct) InsertUser(ctx context.Context, Request *pb.RequestInsert) (*pb.ResponseInsert, error) {

	UserID := UCDataBaseStruct.GenerateUserLastID()
	UserRoleID := UCDataBaseStruct.GenerateUserRoleLastID()
	timeUser := time.Now()
	IUser := dblayer.NewUser(UserID, Request.GetUser().GetUserFname(), Request.GetUser().GetUserLname(),
		Request.GetUser().GetUserMobile(), Request.GetUser().GetUserAddress(), Request.GetUser().GetUserCodeMeli(),
		timeUser, Request.GetUser().GetUserIsActive())
	stateUser := UCDataBaseStruct.InsertUser(IUser)
	IUserRole := dblayer.NewUserRole(UserRoleID, UserID, 2)
	stateUserRole := UCDataBaseStruct.InsertUserRole(IUserRole)

	if stateUser && stateUserRole {
		return &pb.ResponseInsert{
			Status: &pb.Status{
				StatusCode:    pb.StatusCode_Status200,
				StatusMessage: pb.StatusMessage_SUCCESS,
			},
		}, nil
	}
	return &pb.ResponseInsert{
		Status: &pb.Status{
			StatusCode:    pb.StatusCode_Status400,
			StatusMessage: pb.StatusMessage_FAILED,
		},
	}, nil
}
