package utils

import (
	"github.com/google/uuid"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/list/db"
	"time"
)

func newSeedNote(content, attachment string) db.Note {
	return db.Note{
		UserId:     "123",
		NoteId:     uuid.New().String(),
		Content:    content,
		Attachment: attachment,
		CreatedAt:  time.Now().Format(time.RFC3339Nano),
	}
}

func GenerateSeedNotes() []db.Note {
	return []db.Note{
		newSeedNote("pakyu", "betlog.jpg"),
		newSeedNote("tanginamo", "tatlong_itlog.png"),
		newSeedNote("pansit malabon", "minsan.ogg"),
		newSeedNote("kahit pa kingina naman talaga", "you.mp4"),
		newSeedNote("42069 tangninaka hayup ka pakshet", ""),
	}
}
