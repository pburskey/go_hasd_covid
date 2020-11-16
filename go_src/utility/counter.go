package utility

type Counter struct {
	id uint
}

func (c *Counter) next() *Counter {
	c.id = c.id + 1
	return c
}
func (c *Counter) NextId() uint {
	return (*c).next().id
}
