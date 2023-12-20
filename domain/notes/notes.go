package notes

import (
	//	"strings"
	"github.com/google/uuid"
	"time"
	// "golang.org/x/exp/slices"
)

type Note struct {
	CreatedAt time.Time
	id        uuid.UUID
	title     string
	content   string
}

type List struct {
	notes []Note
}

func (notelist *List) Add(title string, content string) {
	notelist.notes = append(notelist.notes, Note{
		CreatedAt: time.Now(),
		content:   content,
		title:     title,
		id:        uuid.New(),
	})
}
