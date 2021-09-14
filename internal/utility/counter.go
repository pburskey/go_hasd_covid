package utility

type Counter struct {
	id int
}

func (c *Counter) next() *Counter {
	c.id = c.id + 1
	return c
}
func (c *Counter) NextId() int {
	return (*c).next().id
}
