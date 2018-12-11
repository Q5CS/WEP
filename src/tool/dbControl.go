package tool

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	//匿名导入数据库驱动
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var marketPlace string
var isMarketPlaceModified bool

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

//CreateConnection 启动数据库连接
func CreateConnection(password string) bool {
	newDB, err := sql.Open("mysql", "app:"+password+"@tcp(localhost:3306)/q5xy?charset=utf8")
	checkErr(err)
	db = newDB
	isMarketPlaceModified = true
	PraseMarketPlace()
	return true
}

//NewUser 新建并记录用户信息（主要是班级）
func NewUser(id string, name string, class string) bool {
	_, err := db.Exec("insert into users (id,name,class) values (?,?,?)", id, name, class)
	if err != nil {
		return false
	}
	return true
}

//HaveUser 检查数据库中是否有该用户
func HaveUser(id string) bool {
	rows, err := db.Query("select count(*) from users where id=?", id)
	if err != nil {
		return false
	}
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return false
		}
	}
	return (count == 1)
}

//HavePermission 检查用户是否拥有某订单的权限
func HavePermission(uID string, role string, orderID string) bool {
	rows, err := db.Query("select count(*) from orders where id=?", orderID)
	if err != nil {
		return false
	}
	var oCount int
	for rows.Next() {
		err = rows.Scan(&oCount)
		if err != nil {
			return false
		}
	}
	result := true
	if oCount != 1 {
		result = false
	} else {
		if role == "0" {
			rows, err = db.Query("select creator from orders where id=?", orderID)
			if err != nil {
				return false
			}
			var cID string
			for rows.Next() {
				err = rows.Scan(&cID)
				if err != nil {
					return false
				}
			}
			if cID != uID {
				result = false
			}
		} else if role == "1" {
			rows, err = db.Query("select matcher from orders where id=?", orderID)
			if err != nil {
				return false
			}
			var mID string
			for rows.Next() {
				err = rows.Scan(&mID)
				if err != nil {
					return false
				}
			}
			if mID != uID {
				result = false
			}
		}
	}
	return result
}

//Create 创建新订单
func Create(id string, item string, amount string, kind string) bool {
	_, err := db.Exec("insert into orders (creator,item,amount,kind,date,status) values (?,?,?,?,now(),0)", id, item, amount, kind)
	if err != nil {
		return false
	}
	isMarketPlaceModified = true
	return true
}

//Match 配对订单
func Match(uID string, orderID string) error {
	rows, err := db.Query("select creator,status from orders where id=?", orderID)
	if err != nil {
		return err
	}
	var creator, status string
	for rows.Next() {
		err = rows.Scan(&creator, &status)
		if err != nil {
			return err
		}
	}
	if creator == uID {
		return errors.New("Selfing")
	} else if status != "0" {
		return errors.New("Invalid Status")
	}
	_, err = db.Exec("update orders set matcher=?,status=1 where id=?", uID, orderID)
	if err != nil {
		return err
	}
	isMarketPlaceModified = true
	return nil
}

//GetInfo 获取另一个用户的班级信息
func GetInfo(uID string, role string, orderID string) (bool, string, string) {
	rows, err := db.Query("select creator,matcher,status from orders where id=?", orderID)
	if err != nil {
		return false, "", ""
	}
	var creator, matcher, status string
	name, class := "", ""
	for rows.Next() {
		err = rows.Scan(&creator, &matcher, &status)
		if err != nil {
			return false, "", ""
		}
	}
	if status == "1" || status == "2" || status == "3" {
		stmt, err := db.Prepare("select name,class from users where id=?")
		if err != nil {
			return false, "", ""
		}

		if role == "0" {
			rows, err := stmt.Query(matcher)
			if err != nil {
				return false, "", ""
			}
			for rows.Next() {
				err = rows.Scan(&name, &class)
				if err != nil {
					return false, "", ""
				}
			}
		} else if role == "1" {
			rows, err := stmt.Query(creator)
			if err != nil {
				return false, "", ""
			}
			for rows.Next() {
				err = rows.Scan(&name, &class)
				if err != nil {
					return false, "", ""
				}
			}
		}
	}
	return true, name, class
}

//Delete 删除订单
func Delete(orderID string) bool {
	_, err := db.Exec("delete from orders where id=?", orderID)
	if err != nil {
		return false
	}
	return true
}

//Reject 拒绝配对请求
func Reject(orderID string) (bool, string, string, string, string) {
	rows, err := db.Query("select creator,item,amount,kind,status from orders where id=?", orderID)
	if err != nil {
		return false, "", "", "", ""
	}
	var creator, item, amount, kind, status string
	for rows.Next() {
		err = rows.Scan(&creator, &item, &amount, &kind, &status)
		if err != nil {
			return false, "", "", "", ""
		}
	}
	if status == "1" {
		_, err = db.Exec("update orders set status=\"3\" where id=?", orderID)
		if err != nil {
			return false, "", "", "", ""
		}
	}
	isMarketPlaceModified = true
	return true, creator, item, amount, kind
}

//Cancel 配对后取消配对请求
func Cancel(orderID string) (bool, string, string, string, string) {
	rows, err := db.Query("select creator,item,amount,kind,status from orders where id=?", orderID)
	if err != nil {
		return false, "", "", "", ""
	}
	var creator, item, amount, kind, status string
	for rows.Next() {
		err = rows.Scan(&creator, &item, &amount, &kind, &status)
		if err != nil {
			return false, "", "", "", ""
		}
	}
	if status == "1" {
		_, err = db.Exec("update orders set status=\"3\" where id=?", orderID)
		if err != nil {
			return false, "", "", "", ""
		}
	}
	isMarketPlaceModified = true
	return true, creator, item, amount, kind
}

//Confirm 确认已收到相应物品
func Confirm(orderID string) bool {
	rows, err := db.Query("select status from orders where id=?", orderID)
	if err != nil {
		return false
	}
	var status string
	for rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			return false
		}
	}
	if status == "1" {
		_, err = db.Exec("update orders set status=\"2\" where id=?", orderID)
		if err != nil {
			return false
		}
	}
	return true
}

//PraseTable 返回"/dashboard"页面的表格
func PraseTable(id string) (bool, string) {
	rows, err := db.Query("select count(*) from orders where creator=?", id)
	if err != nil {
		return false, ""
	}
	var cCount int
	for rows.Next() {
		err = rows.Scan(&cCount)
		if err != nil {
			return false, ""
		}
	}
	rows, err = db.Query("select count(*) from orders where matcher=?", id)
	if err != nil {
		return false, ""
	}
	var mCount int
	for rows.Next() {
		err = rows.Scan(&mCount)
		if err != nil {
			return false, ""
		}
	}

	result := ""
	if cCount == 0 && mCount == 0 {
		result = "Empty Set"
	} else {
		if cCount != 0 {
			rows, err := db.Query("select id,item,amount,kind,date,status from orders where creator=?", id)
			if err != nil {
				return false, ""
			}
			for rows.Next() {
				var (
					id, item, amount, kind, date, status string
				)
				err := rows.Scan(&id, &item, &amount, &kind, &date, &status)
				if err != nil {
					return false, ""
				}
				result += id + "|" + item + "|" + amount + "|" + kind + "|" + date + "|" + status + "|0/"
			}
		}
		result += "||"
		if mCount != 0 {
			rows, err := db.Query("select id,item,amount,kind,date,status from orders where matcher=?", id)
			if err != nil {
				return false, ""
			}
			for rows.Next() {
				var (
					id, item, amount, kind, date, status string
				)
				err := rows.Scan(&id, &item, &amount, &kind, &date, &status)
				if err != nil {
					return false, ""
				}
				result += id + "|" + item + "|" + amount + "|" + kind + "|" + date + "|" + status + "|1/"
			}
		}
	}
	return true, result
}

//PraseMarketPlace 返回/marketPlace页面的表格
func PraseMarketPlace() (bool, string) {
	if isMarketPlaceModified {
		stmt, err := db.Prepare("select id,creator,amount,kind,date from orders where item=? and status=0")
		if err != nil {
			return false, ""
		}
		marketPlace = ""
		for i := 0; i < 5; i++ {
			rows, err := stmt.Query(i)
			if err != nil {
				return false, ""
			}
			/* if !rows.Next() {
				marketPlace += "-||"
				continue
			} */
			flag := false
			for rows.Next() {
				flag = true
				var id, creator, amount, kind, date, class string
				err = rows.Scan(&id, &creator, &amount, &kind, &date)
				if err != nil {
					return false, ""
				}
				user, err := db.Query("select class from users where id=?", creator)
				if err != nil {
					return false, ""
				}
				for user.Next() {
					user.Scan(&class)
				}
				marketPlace += id + "|" + amount + "|" + kind + "|" + date + "|" + class + "/"
			}
			if !flag {
				marketPlace += "-"
			}
			marketPlace += "||"
		}
	}
	return true, marketPlace
}

//Log 自定义的日志函数
func Log(action string, user string, ip string, status string) {
	_, err := db.Exec("insert into logs (time,action,user,ip,status) values (now(),?,?,?,?)", action, user, ip, status)
	if status == "fatal" {
		fmt.Printf(`Fatal error when user "%s" from IP "%s" is performing action "%s"`, user, ip, action)
		os.Exit(1)
	} else if err != nil {
		fmt.Println(`Fail to save logs`)
		os.Exit(1)
	}
}
