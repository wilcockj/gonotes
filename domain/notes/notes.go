package notes

import (
	//	"strings"
	"github.com/google/uuid"
	"time"
	// "golang.org/x/exp/slices"
)

type Note struct {
	CreatedAt time.Time
	NoteUuid  string
	UserId    string
	Title     string
	Content   string
}

type List struct {
	Notes []Note
}

func (notelist *List) Add(user_id string, title string, content string) {
	notelist.Notes = append(notelist.Notes, Note{
		CreatedAt: time.Now(),
		Content:   content,
		Title:     title,
		NoteUuid:  uuid.New().String(),
		UserId:    user_id,
	})
}
