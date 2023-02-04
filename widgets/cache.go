package widgets

type WidgetsCache struct {
	m map[string]Widget
}

func New() *WidgetsCache {
	c := WidgetsCache{
		m: make(map[string]Widget, 100),
	}
	return &c
}

func (c *WidgetsCache) Add(key string, val Widget) bool {
	_, ok := c.Get(key)
	if ok {
		return false
	}
	c.m[key] = val
	return true
}
func (c *WidgetsCache) Map() map[string]Widget {
	return c.m
}
func (c *WidgetsCache) Get(key string) (Widget, bool) {

	v, ok := c.m[key]
	return v, ok
}
