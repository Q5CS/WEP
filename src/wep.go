package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type user struct {
	SessionID string
}

var users map[string]user

var (
	p string
	s string
)

func init() {
	flag.StringVar(&p, "p", "", "Your `mysql_password`")
	flag.StringVar(&s, "s", "", "Your `client_secret` from open.qz5z.ren")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Workbook Exchange Platform
Usage: wep -p mysql_password -s client_secret

Options:
`)
		flag.PrintDefaults()
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	flag.Parse()
	if p == "" || s == "" {
		flag.Usage()
		os.Exit(1)
	}

	users = make(map[string]user)
	err := createConnection(p)
	if err != nil {
		log.Fatalln(err)
	}
	//log("Start", "sys", "localhost", "Succ")

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./front/css/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./front/fonts/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./front/js/"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./front/img/"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/auth_callback", authCallback)
	http.HandleFunc("/dashboard", dashboard)
	http.HandleFunc("/marketPlace", marketPlace)
	http.HandleFunc("/create", create)

	http.HandleFunc("/handlers/dashboard", handleDashboard)
	http.HandleFunc("/handlers/marketPlace", handleMarketPlace)
	http.HandleFunc("/handlers/oppositeInfo", handleOppositeInfo)
	http.HandleFunc("/handlers/create", handleCreate)
	http.HandleFunc("/handlers/match", handleMatch)
	http.HandleFunc("/handlers/delete", handleDelete)
	http.HandleFunc("/handlers/reject", handleReject)
	http.HandleFunc("/handlers/cancel", handleCancel)
	http.HandleFunc("/handlers/confirm", handleConfirm)
	http.HandleFunc("/handlers/auth_callback", handleAuth)
	http.HandleFunc("/handlers/exit", handleExit)

	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./front/index.html")
	t.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./front/login.html")
	t.Execute(w, nil)
}

func authCallback(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./front/auth_callback.html")
	t.Execute(w, nil)
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./front/dashboard.html")
	t.Execute(w, nil)
}

func marketPlace(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./front/marketPlace.html")
	t.Execute(w, nil)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./front/create.html")
	t.Execute(w, nil)
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	/* cUID, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("Server Failure"))
		log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		w.Write([]byte("Server Failure"))
		log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	} */
	var data struct {
		UID       string `json:"uid"`
		SessionID string `json:"sessionID"`
	}
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &data)
	if users[data.UID].SessionID == data.SessionID {
		err, original := praseTable(data.UID)
		if err != nil {
			w.Write([]byte("Server Failure"))
			log.Panicln(err)
			return
			//log("PraseTable", "Sys", "localhost", "Database Error")
		}
		dashBoardData := base64.StdEncoding.EncodeToString([]byte(original))
		w.Write([]byte(dashBoardData))
		//ip := r.Header.Get("X-Real-Ip")
		//log("DashBoard", data.UID, ip, "Succ")
	} else {
		//ip := r.Header.Get("X-Real-Ip")
		w.Write([]byte("Unauthorized"))
		//log("Dashboard", "Unknown", ip, "Validate Fail")
	}
}

func handleMarketPlace(w http.ResponseWriter, r *http.Request) {
	err, original := praseMarketPlace()
	if err != nil {
		log.Panicln(err)
		return
	}
	data := base64.StdEncoding.EncodeToString([]byte(original))
	w.Write([]byte(data))
}

func handleOppositeInfo(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadCookie", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadSession", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	var data struct {
		Role    string `json:"role"`
		OrderID string `json:"orderID"`
	}
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &data)
	var result struct {
		Status string `json:"status"`
		Name   string `json:"name"`
		Class  string `json:"class"`
	}
	if users[uID].SessionID == sessionID {
		err, havePermission := havePermission(uID, data.Role, data.OrderID)
		if err != nil {
			log.Panicln(err)
		}
		if havePermission {
			err, result.Name, result.Class = getInfo(uID, data.Role, data.OrderID)
			if err != nil {
				result.Status = "Server Failure"
				log.Panicln(err)
				return
				//log("GetInfo", "Sys", "localhost", "Database Error")
			}
			//ip := r.Header.Get("X-Real-Ip")
			//log("GetOppositeInfo", uID, ip, "Succ")
			result.Status = "Success"
		} else {
			//ip := r.Header.Get("X-Real-Ip")
			//log("GetOppositeInfo", "Unknown", ip, "Validate Fail")
			result.Status = "Unauthorized"
		}
	} else {
		//ip := r.Header.Get("X-Real-Ip")
		//log("GetOppositeInfo", "Unknown", ip, "Validate Fail")
		result.Status = "Unauthorized"
	}
	b, err = json.Marshal(result)
	w.Write(b)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadCookie", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadSession", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	var data struct {
		Item   string `json:"item"`
		Amount string `json:"amount"`
		Kind   string `json:"kind"`
	}
	b, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &data)
	if users[uID].SessionID == sessionID {
		err = createNew(uID, data.Item, data.Amount, data.Kind)
		if err != nil {
			//log("Create", "Sys", "localhost", "Database Error")
			w.Write([]byte("Server Failure"))
			log.Panicln(err)
			return
		}
		//ip := r.Header.Get("X-Real-Ip")
		//log("HandleCreate", uID, ip, "Succ")
		w.Write([]byte("Success"))
	} else {
		//ip := r.Header.Get("X-Real-Ip")
		//log("HandleCreate", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleMatch(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadCookie", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadSession", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, _ := ioutil.ReadAll(r.Body)
	if users[uID].SessionID == sessionID {
		err = match(uID, string(orderID))
		if err != nil {
			if err.Error() == "Selfing" || err.Error() == "Invalid Status" {
				w.Write([]byte(err.Error()))
			} else {
				//log("Match", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
				log.Panicln(err)
				return
			}
		} else {
			//ip := r.Header.Get("X-Real-Ip")
			//log("HandleMatch", uID, ip, "Succ")
			w.Write([]byte("Success"))
		}
	} else {
		//ip := r.Header.Get("X-Real-Ip")
		//log("HandleMatch", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadCookie", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadSession", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, _ := ioutil.ReadAll(r.Body)
	if users[uID].SessionID == sessionID {
		err, havePermission := havePermission(uID, "0", string(orderID))
		if err != nil {
			log.Panicln(err)
			return
		}
		if havePermission {
			err := deleteOrder(string(orderID))
			if err != nil {
				//log("Delete", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
				log.Panicln(err)
			} else {
				//ip := r.Header.Get("X-Real-Ip")
				//log("HandleDelete", uID, ip, "Succ")
				w.Write([]byte("Success"))
			}
		} else {
			//ip := r.Header.Get("X-Real-Ip")
			//log("HandleDelete", uID, ip, "Permission Denied")
			w.Write([]byte("Unauthorized"))
		}
	} else {
		//ip := r.Header.Get("X-Real-Ip")
		//log("HandleDelete", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleReject(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadCookie", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadSession", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, _ := ioutil.ReadAll(r.Body)
	if users[uID].SessionID == sessionID {
		err, havePermission := havePermission(uID, "0", string(orderID))
		if err != nil {
			//log("Reject", "Sys", "localhost", "Database Error")
			w.Write([]byte("Server Failure"))
			log.Panicln(err)
			return
		}
		if havePermission {
			err, creator, item, amount, kind := reject(string(orderID))
			if err != nil {
				//log("Reject", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
				log.Panicln(err)
				return
			}
			err = createNew(creator, item, amount, kind)
			if err != nil {
				//log("Create After Reject", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
				log.Panicln(err)
				return
			}
			//ip := r.Header.Get("X-Real-Ip")
			//log("HandleReject", uID, ip, "Succ")
			w.Write([]byte("Success"))
		} else {
			//ip := r.Header.Get("X-Real-Ip")
			//log("HandleReject", uID, ip, "Permission Denied")
			w.Write([]byte("Unauthorized"))
		}
	} else {
		//ip := r.Header.Get("X-Real-Ip")
		//log("HandleReject", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleCancel(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadCookie", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadSession", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, _ := ioutil.ReadAll(r.Body)
	if users[uID].SessionID == sessionID {
		err, havePermission := havePermission(uID, "1", string(orderID))
		if err != nil {
			//log("Reject", "Sys", "localhost", "Database Error")
			w.Write([]byte("Server Failure"))
			log.Panicln(err)
			return
		}
		if havePermission {
			err, creator, item, amount, kind := cancel(string(orderID))
			if err != nil {
				//log("Cancel", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
				log.Panicln(err)
				return
			}
			err = createNew(creator, item, amount, kind)
			if err != nil {
				//log("Create After Cancel", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
				log.Panicln(err)
				return
			}
			//ip := r.Header.Get("X-Real-Ip")
			//log("HandleCancel", uID, ip, "Succ")
			w.Write([]byte("Success"))
		} else {
			//ip := r.Header.Get("X-Real-Ip")
			//log("HandleCancel", uID, ip, "Permission Denied")
			w.Write([]byte("Unauthorized"))
		}
	} else {
		//ip := r.Header.Get("X-Real-Ip")
		//log("HandleCancel", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleConfirm(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadCookie", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//log("ReadSession", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, _ := ioutil.ReadAll(r.Body)
	if users[uID].SessionID == sessionID {
		err, havePermission := havePermission(uID, "0", string(orderID))
		if err != nil {
			//log("Reject", "Sys", "localhost", "Database Error")
			w.Write([]byte("Server Failure"))
			log.Panicln(err)
			return
		}
		if havePermission {
			err := confirm(string(orderID))
			if err != nil {
				//log("Confirm", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
				log.Panicln(err)
				return
			}
			//ip := r.Header.Get("X-Real-Ip")
			//log("Confirm", uID, ip, "Succ")
			w.Write([]byte("Success"))
		}
		//ip := r.Header.Get("X-Real-Ip")
		//log("HandleConfirm", uID, ip, "Permission Denied")
		w.Write([]byte("Unauthorized"))
	} else {
		//ip := r.Header.Get("X-Real-Ip")
		//log("HandleConfirm", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	code := r.PostFormValue("code")
	form := fmt.Sprintf("client_id=wep&client_secret=%s&grant_type=authorization_code&code=%s&redirect_uri=https://wep.qz5z.ren/auth_callback&scope=", s, code)
	resp, err := http.Post("https://open.qz5z.ren/oauth2/authorize/token", "application/x-www-form-urlencoded", strings.NewReader(form))
	if err != nil {
		//log("HandleAuth", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
	}

	var auth struct {
		AccessToken string `json:"access_token"`
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &auth)

	form = fmt.Sprintf("access_token=%s&scope=", auth.AccessToken)
	resp, err = http.Post("https://open.qz5z.ren/oauth2/api/getUserData", "application/x-www-form-urlencoded", strings.NewReader(form))
	if err != nil {
		//log("HandleAuth", "Sys", "localhost", "Controller Error")
		log.Panicln(err)
	}

	var data struct {
		UID   string `json:"uid"`
		Name  string `json:"name"`
		Grade string `json:"grade"`
		Class string `json:"class"`
	}
	body, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &data)

	var result struct {
		Status    string `json:"status"`
		UID       string `json:"uID"`
		SessionID string `json:"sessionID"`
		Name      string `json:"name"`
	}
	err, haveUser := haveUser(data.UID)
	if err != nil {
		log.Panicln(err)
	}
	if !haveUser {
		err := newUser(data.UID, data.Name, data.Grade+data.Class)
		if err != nil {
			//log("NewUser", "Sys", "localhost", "Database Error")
			result.Status = "fail"
			result.UID = ""
			b, _ := json.Marshal(result)
			w.Write(b)
			log.Panicln(err)
			return
		}
	}

	var tempUser user
	tempUID, _ := strconv.Atoi(data.UID)
	tempUser.SessionID = strconv.Itoa(tempUID*6 + 233)
	users[data.UID] = tempUser
	result.Status, result.UID, result.SessionID, result.Name = "succ", data.UID, tempUser.SessionID, data.Name
	b, _ := json.Marshal(result)
	//ip := r.Header.Get("X-Real-Ip")
	//log("Login", data.UID, ip, "Succ")
	w.Write(b)
}

func handleExit(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	delete(users, string(b))
	//ip := r.Header.Get("X-Real-Ip")
	//log("HandleExit", string(b), ip, "Succ")
}
