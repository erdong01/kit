package core

func (c *Core) beginTransaction() {
	c.Transaction = c.Db.Begin()
}

// 全局事务
func GlobalTransaction() *Core {
	c := New()
	c.Transaction = c.Db.Begin()
	return c
}
