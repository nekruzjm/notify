package config

func (c *config) Get(key string) any { return c.cfg.Get(key) }

func (c *config) GetBool(key string) bool { return c.cfg.GetBool(key) }

func (c *config) GetFloat64(key string) float64 { return c.cfg.GetFloat64(key) }

func (c *config) GetInt(key string) int { return c.cfg.GetInt(key) }

func (c *config) GetString(key string) string { return c.cfg.GetString(key) }

func (c *config) GetStringSlice(key string) []string {
	return c.cfg.GetStringSlice(key)
}
