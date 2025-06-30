package mysql

// MYSQL测试工具示例

// recordStats 记录用户浏览产品信息
//func recordStats(db *sql.DB, userID, productID int64) (err error) {
//	// 开启事务
//	// 操作views和product_viewers两张表
//	tx, err := db.Begin()
//	if err != nil {
//		return
//	}
//
//	defer func() {
//		switch err {
//		case nil:
//			err = tx.Commit()
//		default:
//			tx.Rollback()
//		}
//	}()
//
//	// 更新products表
//	if _, err = tx.Exec("UPDATE products SET views = views + 1"); err != nil {
//		return
//	}
//	// product_viewers表中插入一条数据
//	if _, err = tx.Exec(
//		"INSERT INTO product_viewers (user_id, product_id) VALUES (?, ?)",
//		userID, productID); err != nil {
//		return
//	}
//	return
//}
//
//func example() {
//	// 注意：测试的过程中并不需要真正的连接
//	db, err := sql.Open("mysql", "root@/blog")
//	if err != nil {
//		panic(err)
//	}
//	defer db.Close()
//	// userID为1的用户浏览了productID为5的产品
//	if err = recordStats(db, 1 /*some user id*/, 5 /*some product id*/); err != nil {
//		panic(err)
//	}
//}
//
//// TestShouldUpdateStats sql执行成功的测试用例
//func TestShouldUpdateStats(t *testing.T) {
//	// mock一个*sql.DB对象，不需要连接真实的数据库
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer db.Close()
//
//	// mock执行指定SQL语句时的返回结果
//	mock.ExpectBegin()
//	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
//	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
//	mock.ExpectCommit()
//
//	// 将mock的DB对象传入我们的函数中
//	if err = recordStats(db, 2, 3); err != nil {
//		t.Errorf("error was not expected while updating stats: %s", err)
//	}
//
//	// 确保期望的结果都满足
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//}
//
//// TestShouldRollbackStatUpdatesOnFailure sql执行失败回滚的测试用例
//func TestShouldRollbackStatUpdatesOnFailure(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer db.Close()
//
//	mock.ExpectBegin()
//	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
//	mock.ExpectExec("INSERT INTO product_viewers").
//		WithArgs(2, 3).
//		WillReturnError(fmt.Errorf("some error"))
//	mock.ExpectRollback()
//
//	// now we execute our method
//	if err = recordStats(db, 2, 3); err == nil {
//		t.Errorf("was expecting an error, but there was none")
//	}
//
//	// we make sure that all expectations were met
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//}
