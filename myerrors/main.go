package main

import (
	"database/sql"
	"log"
	"myerrors/dao"

	"github.com/pkg/errors"
)

/*
当数据库操作时，比如dao层中遇到一个 sql.ErrNoRows，最好应该 Wrap这个error，抛给上层。
因为 dao 层不应该处理异常，事务的边界不在DAO层上，而是在调用层上，故需要往上抛。另外，
Wrap方法（非唯一，go1,13以上fmt.Errorf+w% 也可以）就能包装底层错误，增加上下文文本
信息并附加调用栈，避免在多层调用时因为增加错误信息而导致根因的丢失。
*/
func main() {
	var s string = ""
	_, err := dao.GetUserById(s)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			//匹配到该错误时，再做进一步处理
			log.Printf("该用户没有找到: %+v", err)
		} else {
			log.Fatal("unknown error")
		}
	}

}
