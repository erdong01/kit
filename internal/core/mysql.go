package core

func (c *Core) beginTransaction() {
	c.Transaction = c.Db.Begin()
}

// 全局事务 禁止不同微服务之间使用 只能在单个微服务中使用
func GlobalTransaction() *Core {
	c := New()
	c.Transaction = c.Db.Begin()
	return c
}
