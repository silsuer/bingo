package mysql

import "sync"

type Handle func()

// 事务
func (m *Mysql) Transaction(handle Handle) *Mysql  {
	// 首先加锁
   var mu *sync.Mutex
   mu.Lock()
   res,err:=m.Exec("START TRANSACTION")
   m.checkAppendError(err)
   m.Results = append(m.Results,res)  // 开启一个事务
   // 执行事务中的数据
   handle()
   // 提交事务
   res,err =m.Exec("COMMIT")
   m.checkAppendError(err)
   m.Results = append(m.Results,res)
   // 解锁
   mu.Unlock()
   return m
}
