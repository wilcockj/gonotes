package notes

import (
	//	"strings"
	"github.com/google/uuid"
	"time"
	// "golang.org/x/exp/slices"
)

type Note struct {
	CreatedAt time.Time
	Id        uuid.UUID
	Title     string
	Content   string
}

type List struct {
	Notes []Note
}

func (notelist *List) Add(title string, content string) {
	notelist.Notes = append(notelist.Notes, Note{
		CreatedAt: time.Now(),
		Content:   content,
		Title:     title,
		Id:        uuid.New(),
	})
}
