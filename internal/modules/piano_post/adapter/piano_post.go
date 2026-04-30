package adapter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	piano_postv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
)

func ToPianoPost(view *gateway.PianoPostView) *piano_postv1.PianoPost {
	if view == nil || view.Post == nil {
		return nil
	}
	p := view.Post
	pb := &piano_postv1.PianoPost{
		Name:         resourcename.FormatPianoPost(p.PianoID, p.ID),
		VisitTime:    timestamppb.New(p.VisitTime),
		Rating:       int32(p.Rating),
		Body:         p.Body,
		Visibility:   toPbVisibility(p.Visibility),
		CommentCount: p.CommentCount,
		LikeCount:    p.LikeCount,
		ViewerLiked:  view.ViewerLiked,
		CreateTime:   timestamppb.New(p.CreateTime),
		UpdateTime:   timestamppb.New(p.UpdateTime),
		PianoName:    resourcename.FormatPiano(p.PianoID),
	}
	if v := p.AmbientNoise; v != nil {
		x := int32(*v)
		pb.AmbientNoise = &x
	}
	if v := p.FootTraffic; v != nil {
		x := int32(*v)
		pb.FootTraffic = &x
	}
	if v := p.Resonance; v != nil {
		x := int32(*v)
		pb.Resonance = &x
	}
	if v := p.KeyTouchWeight; v != nil {
		x := int32(*v)
		pb.KeyTouchWeight = &x
	}
	if v := p.TuningQuality; v != nil {
		x := int32(*v)
		pb.TuningQuality = &x
	}
	if view.AuthorCustomID != "" {
		pb.Author = resourcename.FormatUser(view.AuthorCustomID)
	}
	pb.AuthorDisplayName = view.AuthorDisplayName
	pb.PianoDisplayName = view.PianoDisplayName
	return pb
}

func toPbVisibility(v entity.PostVisibility) piano_postv1.PostVisibility {
	switch v {
	case entity.PostVisibilityPublic:
		return piano_postv1.PostVisibility_POST_VISIBILITY_PUBLIC
	case entity.PostVisibilityPrivate:
		return piano_postv1.PostVisibility_POST_VISIBILITY_PRIVATE
	}
	return piano_postv1.PostVisibility_POST_VISIBILITY_UNSPECIFIED
}

func fromPbVisibility(v piano_postv1.PostVisibility) entity.PostVisibility {
	switch v {
	case piano_postv1.PostVisibility_POST_VISIBILITY_PRIVATE:
		return entity.PostVisibilityPrivate
	}
	// UNSPECIFIED は public 扱い (UI が指定しない場合の既定)。
	return entity.PostVisibilityPublic
}
