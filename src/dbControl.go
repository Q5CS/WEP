package main

import (
	"database/sql"
	"errors"

	//匿名导入数据库驱动
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var marketPlaceData string
var isMarketPlaceModified bool

func createConnection(password string) error {
	newDB, err := sql.Open("mysql", "app:"+password+"@tcp(localhost:3306)/q5xy?charset=utf8")
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	db = newDB
	isMarketPlaceModified = true
	praseMarketPlace()
	return nil
}

func newUser(id string, name string, class string) error {
	_, err := db.Exec("insert into users (id,name,class) values (?,?,?)", id, name, class)
	return err
}

func haveUser(id string) (error, bool) {
	rows, err := db.Query("select count(*) from users where id=?", id)
	if err != nil {
		return err, false
	}
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return err, false
		}
	}
	return nil, (count != 1)
}

func havePermission(uID string, role string, orderID string) (error, bool) {
	rows, err := db.Query("select count(*) from orders where id=?", orderID)
	if err != nil {
		return err, false
	}
	var oCount int
	for rows.Next() {
		err = rows.Scan(&oCount)
		if err != nil {
			return err, false
		}
	}
	result := true
	if oCount != 1 {
		result = false
	} else {
		if role == "0" {
			rows, err = db.Query("select creator from orders where id=?", orderID)
			if err != nil {
				return err, false
			}
			var cID string
			for rows.Next() {
				err = rows.Scan(&cID)
				if err != nil {
					return err, false
				}
			}
			if cID != uID {
				result = false
			}
		} else if role == "1" {
			rows, err = db.Query("select matcher from orders where id=?", orderID)
			if err != nil {
				return err, false
			}
			var mID string
			for rows.Next() {
				err = rows.Scan(&mID)
				if err != nil {
					return err, false
				}
			}
			if mID != uID {
				result = false
			}
		}
	}
	return nil, result
}

func createNew(id string, item string, amount string, kind string) error {
	_, err := db.Exec("insert into orders (creator,item,amount,kind,date,status) values (?,?,?,?,now(),0)", id, item, amount, kind)
	if err != nil {
		return err
	}
	isMarketPlaceModified = true
	return nil
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

func getInfo(uID string, role string, orderID string) (error, string, string) {
	rows, err := db.Query("select creator,matcher,status from orders where id=?", orderID)
	if err != nil {
		return err, "", ""
	}
	var creator, matcher, status string
	name, class := "", ""
	for rows.Next() {
		err = rows.Scan(&creator, &matcher, &status)
		if err != nil {
			return err, "", ""
		}
	}
	if status == "1" || status == "2" || status == "3" {
		stmt, err := db.Prepare("select name,class from users where id=?")
		if err != nil {
			return err, "", ""
		}

		if role == "0" {
			rows, err := stmt.Query(matcher)
			if err != nil {
				return err, "", ""
			}
			for rows.Next() {
				err = rows.Scan(&name, &class)
				if err != nil {
					return err, "", ""
				}
			}
		} else if role == "1" {
			rows, err := stmt.Query(creator)
			if err != nil {
				return err, "", ""
			}
			for rows.Next() {
				err = rows.Scan(&name, &class)
				if err != nil {
					return err, "", ""
				}
			}
		}
	}
	return nil, name, class
}

func deleteOrder(orderID string) error {
	_, err := db.Exec("delete from orders where id=?", orderID)
	return err
}

func reject(orderID string) (error, string, string, string, string) {
	rows, err := db.Query("select creator,item,amount,kind,status from orders where id=?", orderID)
	if err != nil {
		return err, "", "", "", ""
	}
	var creator, item, amount, kind, status string
	for rows.Next() {
		err = rows.Scan(&creator, &item, &amount, &kind, &status)
		if err != nil {
			return err, "", "", "", ""
		}
	}
	if status == "1" {
		_, err = db.Exec("update orders set status=\"3\" where id=?", orderID)
		if err != nil {
			return err, "", "", "", ""
		}
	}
	isMarketPlaceModified = true
	return nil, creator, item, amount, kind
}

func cancel(orderID string) (error, string, string, string, string) {
	rows, err := db.Query("select creator,item,amount,kind,status from orders where id=?", orderID)
	if err != nil {
		return err, "", "", "", ""
	}
	var creator, item, amount, kind, status string
	for rows.Next() {
		err = rows.Scan(&creator, &item, &amount, &kind, &status)
		if err != nil {
			return err, "", "", "", ""
		}
	}
	if status == "1" {
		_, err = db.Exec("update orders set status=\"3\" where id=?", orderID)
		if err != nil {
			return err, "", "", "", ""
		}
	}
	isMarketPlaceModified = true
	return nil, creator, item, amount, kind
}

func confirm(orderID string) error {
	rows, err := db.Query("select status from orders where id=?", orderID)
	if err != nil {
		return err
	}
	var status string
	for rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			return err
		}
	}
	if status == "1" {
		_, err = db.Exec("update orders set status=\"2\" where id=?", orderID)
		if err != nil {
			return err
		}
	}
	return nil
}

func praseTable(id string) (error, string) {
	rows, err := db.Query("select count(*) from orders where creator=?", id)
	if err != nil {
		return err, ""
	}
	var cCount int
	for rows.Next() {
		err = rows.Scan(&cCount)
		if err != nil {
			return err, ""
		}
	}
	rows, err = db.Query("select count(*) from orders where matcher=?", id)
	if err != nil {
		return err, ""
	}
	var mCount int
	for rows.Next() {
		err = rows.Scan(&mCount)
		if err != nil {
			return err, ""
		}
	}

	result := ""
	if cCount == 0 && mCount == 0 {
		result = "Empty Set"
	} else {
		if cCount != 0 {
			rows, err := db.Query("select id,item,amount,kind,date,status from orders where creator=?", id)
			if err != nil {
				return err, ""
			}
			for rows.Next() {
				var (
					id, item, amount, kind, date, status string
				)
				err := rows.Scan(&id, &item, &amount, &kind, &date, &status)
				if err != nil {
					return err, ""
				}
				result += id + "|" + item + "|" + amount + "|" + kind + "|" + date + "|" + status + "|0/"
			}
		}
		result += "||"
		if mCount != 0 {
			rows, err := db.Query("select id,item,amount,kind,date,status from orders where matcher=?", id)
			if err != nil {
				return err, ""
			}
			for rows.Next() {
				var (
					id, item, amount, kind, date, status string
				)
				err := rows.Scan(&id, &item, &amount, &kind, &date, &status)
				if err != nil {
					return err, ""
				}
				result += id + "|" + item + "|" + amount + "|" + kind + "|" + date + "|" + status + "|1/"
			}
		}
	}
	return nil, result
}

func praseMarketPlace() (error, string) {
	if isMarketPlaceModified {
		stmt, err := db.Prepare("select id,creator,amount,kind,date from orders where item=? and status=0")
		if err != nil {
			return err, ""
		}
		marketPlaceData = ""
		for i := 0; i < 5; i++ {
			rows, err := stmt.Query(i)
			if err != nil {
				return err, ""
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
					return err, ""
				}
				user, err := db.Query("select class from users where id=?", creator)
				if err != nil {
					return err, ""
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
	return nil, marketPlaceData
}

/* func log(action string, user string, ip string, status string) {
	_, err := db.Exec("insert into logs (time,action,user,ip,status) values (now(),?,?,?,?)", action, user, ip, status)
	if status == "fatal" {
		fmt.Printf(`Fatal error when user "%s" from IP "%s" is performing action "%s"`, user, ip, action)
		os.Exit(1)
	} else if err != nil {
		fmt.Println(`Fail to save logs`)
		os.Exit(1)
	}
} */
