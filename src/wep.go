package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
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
}

func main() {
	flag.Parse()
	if p == "" || s == "" {
		flag.Usage()
		os.Exit(1)
	}

	users = make(map[string]user)
	createConnection(p)
	log("Start", "sys", "localhost", "Succ")

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

	err := http.ListenAndServe(":9090", nil)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
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
		succ, original := praseTable(data.UID)
		if !succ {
			w.Write([]byte("Server Failure"))
			log("PraseTable", "Sys", "localhost", "Database Error")
		} else {
			dashBoardData := base64.StdEncoding.EncodeToString([]byte(original))
			w.Write([]byte(dashBoardData))
			ip := r.Header.Get("X-Real-Ip")
			log("DashBoard", data.UID, ip, "Succ")
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		w.Write([]byte("Unauthorized"))
		log("Dashboard", "Unknown", ip, "Validate Fail")
	}
}

func handleMarketPlace(w http.ResponseWriter, r *http.Request) {
	succ, original := praseMarketPlace()
	if !succ {
		log("MarketPlace", "Sys", "localhost", "Database Error")
	}
	data := base64.StdEncoding.EncodeToString([]byte(original))
	w.Write([]byte(data))
}

func handleOppositeInfo(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadSession", "Sys", "localhost", "Controller Error")
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
		if havePermission(uID, data.Role, data.OrderID) {
			var succ bool
			succ, result.Name, result.Class = getInfo(uID, data.Role, data.OrderID)
			if !succ {
				log("GetInfo", "Sys", "localhost", "Database Error")
				result.Status = "Server Failure"
			} else {
				ip := r.Header.Get("X-Real-Ip")
				log("GetOppositeInfo", uID, ip, "Succ")
				result.Status = "Success"
			}
		} else {
			ip := r.Header.Get("X-Real-Ip")
			log("GetOppositeInfo", "Unknown", ip, "Validate Fail")
			result.Status = "Unauthorized"
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		log("GetOppositeInfo", "Unknown", ip, "Validate Fail")
		result.Status = "Unauthorized"
	}
	b, err = json.Marshal(result)
	w.Write(b)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadSession", "Sys", "localhost", "Controller Error")
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
		succ := createNew(uID, data.Item, data.Amount, data.Kind)
		if !succ {
			log("Create", "Sys", "localhost", "Database Error")
			w.Write([]byte("Server Failure"))
		} else {
			ip := r.Header.Get("X-Real-Ip")
			log("HandleCreate", uID, ip, "Succ")
			w.Write([]byte("Success"))
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		log("HandleCreate", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleMatch(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, _ := ioutil.ReadAll(r.Body)
	if users[uID].SessionID == sessionID {
		err = match(uID, string(orderID))
		if err != nil {
			if err.Error() == "Selfing" || err.Error() == "Invalid Status" {
				log("Match", "Sys", "localhost", err.Error())
				w.Write([]byte(err.Error()))
			} else {
				log("Match", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
			}
		} else {
			ip := r.Header.Get("X-Real-Ip")
			log("HandleMatch", uID, ip, "Succ")
			w.Write([]byte("Success"))
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		log("HandleMatch", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, _ := ioutil.ReadAll(r.Body)
	if users[uID].SessionID == sessionID {
		if havePermission(uID, "0", string(orderID)) {
			if !deleteOrder(string(orderID)) {
				log("Delete", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
			} else {
				ip := r.Header.Get("X-Real-Ip")
				log("HandleDelete", uID, ip, "Succ")
				w.Write([]byte("Success"))
			}
		} else {
			ip := r.Header.Get("X-Real-Ip")
			log("HandleDelete", uID, ip, "Permission Denied")
			w.Write([]byte("Unauthorized"))
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		log("HandleDelete", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleReject(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, _ := ioutil.ReadAll(r.Body)
	if users[uID].SessionID == sessionID {
		if havePermission(uID, "0", string(orderID)) {
			succ, creator, item, amount, kind := reject(string(orderID))
			if !succ {
				log("Reject", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
				return
			}
			succ = createNew(creator, item, amount, kind)
			if !succ {
				log("Create After Reject", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
			} else {
				ip := r.Header.Get("X-Real-Ip")
				log("HandleReject", uID, ip, "Succ")
				w.Write([]byte("Success"))
			}
		} else {
			ip := r.Header.Get("X-Real-Ip")
			log("HandleReject", uID, ip, "Permission Denied")
			w.Write([]byte("Unauthorized"))
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		log("HandleReject", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleCancel(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, _ := ioutil.ReadAll(r.Body)
	if users[uID].SessionID == sessionID {
		if havePermission(uID, "1", string(orderID)) {
			succ, creator, item, amount, kind := cancel(string(orderID))
			if !succ {
				log("Cancel", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
				return
			}
			succ = createNew(creator, item, amount, kind)
			if !succ {
				log("Create After Cancel", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
			} else {
				ip := r.Header.Get("X-Real-Ip")
				log("HandleCancel", uID, ip, "Succ")
				w.Write([]byte("Success"))
			}
		} else {
			ip := r.Header.Get("X-Real-Ip")
			log("HandleCancel", uID, ip, "Permission Denied")
			w.Write([]byte("Unauthorized"))
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		log("HandleCancel", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleConfirm(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, _ := ioutil.ReadAll(r.Body)
	if users[uID].SessionID == sessionID {
		if havePermission(uID, "0", string(orderID)) {
			succ := confirm(string(orderID))
			if !succ {
				log("Confirm", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
			} else {
				ip := r.Header.Get("X-Real-Ip")
				log("Confirm", uID, ip, "Succ")
				w.Write([]byte("Success"))
			}
		} else {
			ip := r.Header.Get("X-Real-Ip")
			log("HandleConfirm", uID, ip, "Permission Denied")
			w.Write([]byte("Unauthorized"))
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		log("HandleConfirm", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	code := r.PostFormValue("code")
	form := fmt.Sprintf("client_id=wep&client_secret=%s&grant_type=authorization_code&code=%s&redirect_uri=https://wep.qz5z.ren/auth_callback&scope=", s, code)
	resp, err := http.Post("https://open.qz5z.ren/oauth2/authorize/token", "application/x-www-form-urlencoded", strings.NewReader(form))
	if err != nil {
		log("HandleAuth", "Sys", "localhost", "Controller Error")
	}

	var auth struct {
		AccessToken string `json:"access_token"`
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &auth)

	form = fmt.Sprintf("access_token=%s&scope=", auth.AccessToken)
	resp, err = http.Post("https://open.qz5z.ren/oauth2/api/getUserData", "application/x-www-form-urlencoded", strings.NewReader(form))
	if err != nil {
		log("HandleAuth", "Sys", "localhost", "Controller Error")
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
	if !haveUser(data.UID) {
		succ := newUser(data.UID, data.Name, data.Grade+data.Class)
		if !succ {
			log("NewUser", "Sys", "localhost", "Database Error")
			result.Status = "fail"
			result.UID = ""
			b, _ := json.Marshal(result)
			w.Write(b)
			return
		}
	}

	var tempUser user
	tempUID, _ := strconv.Atoi(data.UID)
	tempUser.SessionID = strconv.Itoa(tempUID*6 + 233)
	users[data.UID] = tempUser
	result.Status, result.UID, result.SessionID, result.Name = "succ", data.UID, tempUser.SessionID, data.Name
	b, _ := json.Marshal(result)
	ip := r.Header.Get("X-Real-Ip")
	log("Login", data.UID, ip, "Succ")
	w.Write(b)
}

func handleExit(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	delete(users, string(b))
	ip := r.Header.Get("X-Real-Ip")
	log("HandleExit", string(b), ip, "Succ")
}
