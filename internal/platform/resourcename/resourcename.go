// Package resourcename translates between internal identifiers and
// AIP-122-style resource names ("collection/id" hierarchical strings).
package resourcename

import (
	"errors"
	"strings"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

const (
	usersCollection         = "users"
	pianosCollection        = "pianos"
	postsCollection         = "posts"
	commentsCollection      = "comments"
	photosCollection        = "photos"
	editsCollection         = "edits"
	notificationsCollection = "notifications"
	reportsCollection       = "reports"
)

// User: users/{custom_id}

func FormatUser(customID string) string {
	return usersCollection + "/" + customID
}

func ParseUser(name string) (string, error) {
	segs, err := split(name)
	if err != nil {
		return "", err
	}
	if len(segs) != 2 || segs[0] != usersCollection || segs[1] == "" {
		return "", errors.New("invalid user resource name: " + name)
	}
	return segs[1], nil
}

// Piano: pianos/{piano_id}

func FormatPiano(id ulid.ULID) string {
	return pianosCollection + "/" + id.String()
}

func ParsePiano(name string) (ulid.ULID, error) {
	segs, err := split(name)
	if err != nil {
		return ulid.ULID{}, err
	}
	if len(segs) != 2 || segs[0] != pianosCollection || segs[1] == "" {
		return ulid.ULID{}, errors.New("invalid piano resource name: " + name)
	}
	return ulid.Parse(segs[1])
}

// PianoPhoto: pianos/{piano_id}/photos/{photo_id}

func FormatPianoPhoto(pianoID, photoID ulid.ULID) string {
	return FormatPiano(pianoID) + "/" + photosCollection + "/" + photoID.String()
}

func ParsePianoPhoto(name string) (pianoID, photoID ulid.ULID, err error) {
	segs, err := split(name)
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	if len(segs) != 4 ||
		segs[0] != pianosCollection ||
		segs[2] != photosCollection ||
		segs[1] == "" || segs[3] == "" {
		return ulid.ULID{}, ulid.ULID{}, errors.New("invalid piano photo resource name: " + name)
	}
	pianoID, err = ulid.Parse(segs[1])
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	photoID, err = ulid.Parse(segs[3])
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	return pianoID, photoID, nil
}

// PianoComment: pianos/{piano_id}/comments/{comment_id}

func FormatPianoComment(pianoID, commentID ulid.ULID) string {
	return FormatPiano(pianoID) + "/" + commentsCollection + "/" + commentID.String()
}

func ParsePianoComment(name string) (pianoID, commentID ulid.ULID, err error) {
	segs, err := split(name)
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	if len(segs) != 4 ||
		segs[0] != pianosCollection ||
		segs[2] != commentsCollection ||
		segs[1] == "" || segs[3] == "" {
		return ulid.ULID{}, ulid.ULID{}, errors.New("invalid piano comment resource name: " + name)
	}
	pianoID, err = ulid.Parse(segs[1])
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	commentID, err = ulid.Parse(segs[3])
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	return pianoID, commentID, nil
}

// PianoEdit: pianos/{piano_id}/edits/{edit_id}

func FormatPianoEdit(pianoID, editID ulid.ULID) string {
	return FormatPiano(pianoID) + "/" + editsCollection + "/" + editID.String()
}

func ParsePianoEdit(name string) (pianoID, editID ulid.ULID, err error) {
	segs, err := split(name)
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	if len(segs) != 4 ||
		segs[0] != pianosCollection ||
		segs[2] != editsCollection ||
		segs[1] == "" || segs[3] == "" {
		return ulid.ULID{}, ulid.ULID{}, errors.New("invalid piano edit resource name: " + name)
	}
	pianoID, err = ulid.Parse(segs[1])
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	editID, err = ulid.Parse(segs[3])
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	return pianoID, editID, nil
}

// PianoPost: pianos/{piano_id}/posts/{post_id}

func FormatPianoPost(pianoID, postID ulid.ULID) string {
	return FormatPiano(pianoID) + "/" + postsCollection + "/" + postID.String()
}

func ParsePianoPost(name string) (pianoID, postID ulid.ULID, err error) {
	segs, err := split(name)
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	if len(segs) != 4 ||
		segs[0] != pianosCollection ||
		segs[2] != postsCollection ||
		segs[1] == "" || segs[3] == "" {
		return ulid.ULID{}, ulid.ULID{}, errors.New("invalid piano post resource name: " + name)
	}
	pianoID, err = ulid.Parse(segs[1])
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	postID, err = ulid.Parse(segs[3])
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, err
	}
	return pianoID, postID, nil
}

// PianoPostComment: pianos/{piano_id}/posts/{post_id}/comments/{comment_id}

func FormatPianoPostComment(pianoID, postID, commentID ulid.ULID) string {
	return FormatPianoPost(pianoID, postID) + "/" + commentsCollection + "/" + commentID.String()
}

func ParsePianoPostComment(name string) (pianoID, postID, commentID ulid.ULID, err error) {
	segs, err := split(name)
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, ulid.ULID{}, err
	}
	if len(segs) != 6 ||
		segs[0] != pianosCollection ||
		segs[2] != postsCollection ||
		segs[4] != commentsCollection ||
		segs[1] == "" || segs[3] == "" || segs[5] == "" {
		return ulid.ULID{}, ulid.ULID{}, ulid.ULID{}, errors.New("invalid piano post comment resource name: " + name)
	}
	pianoID, err = ulid.Parse(segs[1])
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, ulid.ULID{}, err
	}
	postID, err = ulid.Parse(segs[3])
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, ulid.ULID{}, err
	}
	commentID, err = ulid.Parse(segs[5])
	if err != nil {
		return ulid.ULID{}, ulid.ULID{}, ulid.ULID{}, err
	}
	return pianoID, postID, commentID, nil
}

// Notification: users/{custom_id}/notifications/{notification_id}

func FormatNotification(userCustomID string, notificationID ulid.ULID) string {
	return FormatUser(userCustomID) + "/" + notificationsCollection + "/" + notificationID.String()
}

func ParseNotification(name string) (userCustomID string, notificationID ulid.ULID, err error) {
	segs, err := split(name)
	if err != nil {
		return "", ulid.ULID{}, err
	}
	if len(segs) != 4 ||
		segs[0] != usersCollection ||
		segs[2] != notificationsCollection ||
		segs[1] == "" || segs[3] == "" {
		return "", ulid.ULID{}, errors.New("invalid notification resource name: " + name)
	}
	notificationID, err = ulid.Parse(segs[3])
	if err != nil {
		return "", ulid.ULID{}, err
	}
	return segs[1], notificationID, nil
}

// Report: reports/{report_id}

func FormatReport(id ulid.ULID) string {
	return reportsCollection + "/" + id.String()
}

func ParseReport(name string) (ulid.ULID, error) {
	segs, err := split(name)
	if err != nil {
		return ulid.ULID{}, err
	}
	if len(segs) != 2 || segs[0] != reportsCollection || segs[1] == "" {
		return ulid.ULID{}, errors.New("invalid report resource name: " + name)
	}
	return ulid.Parse(segs[1])
}

func split(name string) ([]string, error) {
	if name == "" {
		return nil, errors.New("empty resource name")
	}
	return strings.Split(name, "/"), nil
}
