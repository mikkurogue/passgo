package pkg

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/labstack/gommon/log"
)

type Searchable interface {
	GetRows() []table.Row
}

func Search(t Searchable, inp string) {
	rows := t.GetRows()

	log.Debug(rows)
}
