package models

import (
	"github.com/uptrace/bun"
)

type Note struct {
	bun.BaseModel `bun:"table:notes"`
	ID            int64    `bun:",pk,autoincrement"`
	UUID          string   `bun:",notnull,type:uuid,default:gen_random_uuid()"`
	Title         string   `bun:",notnull"`
	Contents      string   `bun:",notnull"`
	CreatorID     int64    `bun:",notnull"`
	Type          NoteType `bun:"type:varchar(10),notnull"`
}

type NoteType string

const (
	NoteTypeGuild NoteType = "GUILD"
	NoteTypeUser  NoteType = "USER"
)
