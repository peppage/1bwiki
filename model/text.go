package model

type Text struct {
	ID   int64
	Text string
}

func (t Text) Verify() error {
	if t.Text == "" {
		return logger.Error("Empty text")
	}
	return nil
}
