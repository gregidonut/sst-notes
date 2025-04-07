package utils

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/list/db"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/utils"
	"time"
)

func GenerateSeedNotes(event events.APIGatewayProxyRequest) ([]db.Note, error) {
	payload := []db.Note{}

	userId, err := utils.GetUserId(event)
	if err != nil {
		return nil, err
	}

	for _, v := range []struct {
		content     string
		attachement string
	}{

		{"pakyu", "betlog.jpg"},
		{"tanginamo", "tatlong_itlog.png"},
		{"pansit malabon", "minsan.ogg"},
		{"kahit pa kingina naman talaga", "you.mp4"},
		{"42069 tangninaka hayup ka pakshet", ""},
	} {
		payload = append(payload, db.Note{
			UserId:     userId,
			NoteId:     uuid.New().String(),
			Content:    v.content,
			Attachment: v.attachement,
			CreatedAt:  time.Now().Format(time.RFC3339Nano),
		})
	}
	return payload, nil
}
