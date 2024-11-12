package pkg

import (
	"github.com/atotto/clipboard"
)

type Copier interface {
	Copy(value string) error
	Histroy() []string
}

type ClipboardCopier struct {
	history []string
}

func (c *ClipboardCopier) Copy(value string) error {
	if err := clipboard.WriteAll(value); err != nil {
		return err
	}

	c.history = append(c.history, value)
	return nil
}

func (c *ClipboardCopier) History() []string {
	return c.history
}
