package essence

type Closable interface {
	Close() error
}

func Close(c interface{}) error {
	if c == nil {
		return nil
	}

	if closer, ok := c.(Closable); ok {
		return closer.Close()
	}

	return nil
}
