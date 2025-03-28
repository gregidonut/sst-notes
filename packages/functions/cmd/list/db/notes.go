package db

type Note struct {
	Attachment string `json:"attachment" dynamodbav:"attachment"`
	Content    string `json:"content" dynamodbav:"content"`
	CreatedAt  string `json:"createdAt" dynamodbav:"createdAt"`
	NoteId     string `json:"noteId" dynamodbav:"noteId"`
	UserId     string `json:"userId" dynamodbav:"userId"`
}
