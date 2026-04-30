package adapter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "github.com/reverie-jp/piamap/internal/gen/pb/user/v1"
	"github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
)

func ToUser(view *gateway.UserView) *userv1.User {
	if view == nil || view.User == nil {
		return nil
	}
	u := view.User
	pb := &userv1.User{
		Name:        resourcename.FormatUser(u.CustomID),
		CustomId:    u.CustomID,
		DisplayName: u.DisplayName,
		Biography:   u.Biography,
		AvatarUrl:   u.AvatarURL,
		Hometown:    u.Hometown,
		PostCount:   u.PostCount,
		EditCount:   u.EditCount,
		CreateTime:  timestamppb.New(u.CreateTime),
		UpdateTime:  timestamppb.New(u.UpdateTime),
		IsMe:        view.IsMe,
	}
	if u.PianoStartedYear != nil {
		v := int32(*u.PianoStartedYear)
		pb.PianoStartedYear = &v
	}
	if u.YearsOfExperience != nil {
		v := int32(*u.YearsOfExperience)
		pb.YearsOfExperience = &v
	}
	if u.CustomIDChangeTime != nil {
		pb.CustomIdChangeTime = timestamppb.New(*u.CustomIDChangeTime)
	}
	return pb
}
