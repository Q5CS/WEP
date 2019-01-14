package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	//匿名导入数据库驱动
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var marketPlaceData string
var isMarketPlaceModified bool

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func createConnection(password string) bool {
	newDB, err := sql.Open("mysql", "app:"+password+"@tcp(localhost:3306)/q5xy?charset=utf8")
	checkErr(err)
	db = newDB
	isMarketPlaceModified = true
	praseMarketPlace()
	return true
}

func newUser(id string, name string, class string) bool {
	_, err := db.Exec("insert into users (id,name,class) values (?,?,?)", id, name, class)
	if err != nil {
		return false
	}
	return true
}

func haveUser(id string) bool {
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

func havePermission(uID string, role string, orderID string) bool {
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

func createNew(id string, item string, amount string, kind string) bool {
	_, err := db.Exec("insert into orders (creator,item,amount,kind,date,status) values (?,?,?,?,now(),0)", id, item, amount, kind)
	if err != nil {
		return false
	}
	isMarketPlaceModified = true
	return true
}

func match(uID string, orderID string) error {
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

func getInfo(uID string, role string, orderID string) (bool, string, string) {
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

func deleteOrder(orderID string) bool {
	_, err := db.Exec("delete from orders where id=?", orderID)
	if err != nil {
		return false
	}
	return true
}

func reject(orderID string) (bool, string, string, string, string) {
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

func cancel(orderID string) (bool, string, string, string, string) {
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

func confirm(orderID string) bool {
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

func praseTable(id string) (bool, string) {
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

func praseMarketPlace() (bool, string) {
	if isMarketPlaceModified {
		stmt, err := db.Prepare("select id,creator,amount,kind,date from orders where item=? and status=0")
		if err != nil {
			return false, ""
		}
		marketPlaceData = ""
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
				marketPlaceData += id + "|" + amount + "|" + kind + "|" + date + "|" + class + "/"
			}
			if !flag {
				marketPlaceData += "-"
			}
			marketPlaceData += "||"
		}
	}
	return true, marketPlaceData
}

func log(action string, user string, ip string, status string) {
	_, err := db.Exec("insert into logs (time,action,user,ip,status) values (now(),?,?,?,?)", action, user, ip, status)
	if status == "fatal" {
		fmt.Printf(`Fatal error when user "%s" from IP "%s" is performing action "%s"`, user, ip, action)
		os.Exit(1)
	} else if err != nil {
		fmt.Println(`Fail to save logs`)
		os.Exit(1)
	}
}
