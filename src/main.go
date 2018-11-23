package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"tool"
)

type user struct {
	SessionID string
}

var users map[string]user

func main() {
	users = make(map[string]user)
	tool.CreateConnection()
	tool.Log("Start", "sys", "localhost", "Succ")

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./front/css/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./front/fonts/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./front/js/"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./front/img/"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/report", report)
	http.HandleFunc("/dashboard", dashboard)
	http.HandleFunc("/marketPlace", marketPlace)
	http.HandleFunc("/create", create)

	http.HandleFunc("/handlers/oppositeInfo", handleOppositeInfo)
	http.HandleFunc("/handlers/create", handleCreate)
	http.HandleFunc("/handlers/match", handleMatch)
	http.HandleFunc("/handlers/delete", handleDelete)
	http.HandleFunc("/handlers/reject", handleReject)
	http.HandleFunc("/handlers/cancel", handleCancel)
	http.HandleFunc("/handlers/confirm", handleConfirm)
	http.HandleFunc("/handlers/login", handleLogin)
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
	t, err := template.ParseFiles("./front/index.html")
	if err != nil {
		fmt.Println(err)
		tool.Log("Index", "Sys", "localhost", "Controller Error")
	}
	t.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./front/login.html")
	if err != nil {
		fmt.Println(err)
		tool.Log("Login", "Sys", "localhost", "Controller Error")
	}
	t.Execute(w, nil)
}

func report(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./front/report.html")
	if err != nil {
		fmt.Println(err)
		tool.Log("Report", "Sys", "localhost", "Controller Error")
	}
	t.Execute(w, nil)
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	if users[uID].SessionID == sessionID {
		t, err := template.ParseFiles("./front/dashboard.html")
		if err != nil {
			fmt.Println(err)
			tool.Log("Dashboard", "Sys", "localhost", "Controller Error")
		}
		succ, original := tool.PraseTable(uID)
		if !succ {
			tool.Log("PraseTable", "Sys", "localhost", "Database Error")
		}
		data := base64.StdEncoding.EncodeToString([]byte(original))
		t.Execute(w, data)
	} else {
		ip := r.Header.Get("X-Real-Ip")
		tool.Log("Dashboard", "Unknown", ip, "Validate Fail")
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	ip := r.Header.Get("X-Real-Ip")
	tool.Log("DashBoard", uID, ip, "Succ")
}

func marketPlace(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./front/marketPlace.html")
	if err != nil {
		fmt.Println(err)
		tool.Log("MarketPlace", "Sys", "localhost", "Controller Error")
	}
	succ, original := tool.PraseMarketPlace()
	if !succ {
		tool.Log("MarketPlace", "Sys", "localhost", "Database Error")
	}
	data := base64.StdEncoding.EncodeToString([]byte(original))
	t.Execute(w, data)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./front/create.html")
	if err != nil {
		tool.Log("Create", "Sys", "localhost", "Controller Error")
	}
	t.Execute(w, nil)
}

func handleOppositeInfo(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	var data struct {
		Role    string `json:"role"`
		OrderID string `json:"orderID"`
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.Log("HandleOppositeInfo", "Sys", "localhost", "Controller Error")
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		tool.Log("HandleOppositeInfo", "Sys", "localhost", "Controller Error")
	}
	var result struct {
		Status string `json:"status"`
		Name   string `json:"name"`
		Class  string `json:"class"`
	}
	if users[uID].SessionID == sessionID {
		if tool.HavePermission(uID, data.Role, data.OrderID) {
			var succ bool
			succ, result.Name, result.Class = tool.GetInfo(uID, data.Role, data.OrderID)
			if !succ {
				tool.Log("GetInfo", "Sys", "localhost", "Database Error")
				result.Status = "Server Failure"
			} else {
				ip := r.Header.Get("X-Real-Ip")
				tool.Log("GetOppositeInfo", uID, ip, "Succ")
				result.Status = "Success"
			}
		} else {
			ip := r.Header.Get("X-Real-Ip")
			tool.Log("GetOppositeInfo", "Unknown", ip, "Validate Fail")
			result.Status = "Unauthorized"
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		tool.Log("GetOppositeInfo", "Unknown", ip, "Validate Fail")
		result.Status = "Unauthorized"
	}
	b, err = json.Marshal(result)
	w.Write(b)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	var data struct {
		Item   string `json:"item"`
		Amount string `json:"amount"`
		Kind   string `json:"kind"`
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.Log("HandleCreate", "Sys", "localhost", "Controller Error")
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		tool.Log("HandleCreate", "Sys", "localhost", "Controller Error")
	}
	if users[uID].SessionID == sessionID {
		succ := tool.Create(uID, data.Item, data.Amount, data.Kind)
		if !succ {
			tool.Log("Create", "Sys", "localhost", "Database Error")
			w.Write([]byte("Server Failure"))
		} else {
			ip := r.Header.Get("X-Real-Ip")
			tool.Log("HandleCreate", uID, ip, "Succ")
			w.Write([]byte("Success"))
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		tool.Log("HandleCreate", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleMatch(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.Log("HandleMatch", "Sys", "localhost", "Controller Error")
	}
	if users[uID].SessionID == sessionID {
		succ := tool.Match(uID, string(orderID))
		if !succ {
			tool.Log("Match", "Sys", "localhost", "Database Error")
			w.Write([]byte("Server Failure"))
		} else {
			ip := r.Header.Get("X-Real-Ip")
			tool.Log("HandleMatch", uID, ip, "Succ")
			w.Write([]byte("Success"))
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		tool.Log("HandleMatch", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.Log("HandleDelete", "Sys", "localhost", "Controller Error")
	}
	if users[uID].SessionID == sessionID {
		if tool.HavePermission(uID, "0", string(orderID)) {
			if !tool.Delete(string(orderID)) {
				tool.Log("Delete", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
			} else {
				ip := r.Header.Get("X-Real-Ip")
				tool.Log("HandleDelete", uID, ip, "Succ")
				w.Write([]byte("Success"))
			}
		} else {
			ip := r.Header.Get("X-Real-Ip")
			tool.Log("HandleDelete", uID, ip, "Permission Denied")
			w.Write([]byte("Unauthorized"))
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		tool.Log("HandleDelete", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleReject(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.Log("HandleReject", "Sys", "localhost", "Controller Error")
	}
	if users[uID].SessionID == sessionID {
		if tool.HavePermission(uID, "0", string(orderID)) {
			succ, creator, item, amount, kind := tool.Reject(string(orderID))
			if !succ {
				tool.Log("Reject", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
				return
			}
			succ = tool.Create(creator, item, amount, kind)
			if !succ {
				tool.Log("Create After Reject", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
			} else {
				ip := r.Header.Get("X-Real-Ip")
				tool.Log("HandleReject", uID, ip, "Succ")
				w.Write([]byte("Success"))
			}
		} else {
			ip := r.Header.Get("X-Real-Ip")
			tool.Log("HandleReject", uID, ip, "Permission Denied")
			w.Write([]byte("Unauthorized"))
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		tool.Log("HandleReject", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleCancel(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.Log("HandleCancel", "Sys", "localhost", "Controller Error")
	}
	if users[uID].SessionID == sessionID {
		if tool.HavePermission(uID, "1", string(orderID)) {
			succ, creator, item, amount, kind := tool.Cancel(string(orderID))
			if !succ {
				tool.Log("Cancel", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
				return
			}
			succ = tool.Create(creator, item, amount, kind)
			if !succ {
				tool.Log("Create After Cancel", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
			} else {
				ip := r.Header.Get("X-Real-Ip")
				tool.Log("HandleCancel", uID, ip, "Succ")
				w.Write([]byte("Success"))
			}
		} else {
			ip := r.Header.Get("X-Real-Ip")
			tool.Log("HandleCancel", uID, ip, "Permission Denied")
			w.Write([]byte("Unauthorized"))
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		tool.Log("HandleCancel", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleConfirm(w http.ResponseWriter, r *http.Request) {
	cUID, err := r.Cookie("uid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadCookie", "Sys", "localhost", "Controller Error")
		return
	}
	cSessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		tool.Log("ReadSession", "Sys", "localhost", "Controller Error")
		return
	}
	uID, sessionID := cUID.Value, cSessionID.Value
	orderID, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.Log("HandleConfirm", "Sys", "localhost", "Controller Error")
	}
	if users[uID].SessionID == sessionID {
		if tool.HavePermission(uID, "0", string(orderID)) {
			succ := tool.Confirm(string(orderID))
			if !succ {
				tool.Log("Confirm", "Sys", "localhost", "Database Error")
				w.Write([]byte("Server Failure"))
			} else {
				ip := r.Header.Get("X-Real-Ip")
				tool.Log("Confirm", uID, ip, "Succ")
				w.Write([]byte("Success"))
			}
		} else {
			ip := r.Header.Get("X-Real-Ip")
			tool.Log("HandleConfirm", uID, ip, "Permission Denied")
			w.Write([]byte("Unauthorized"))
		}
	} else {
		ip := r.Header.Get("X-Real-Ip")
		tool.Log("HandleConfirm", "Unknown", ip, "Validate Fail")
		w.Write([]byte("Unauthorized"))
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.Log("HandleLogin", "Sys", "localhost", "Controller Error")
	}
	var data struct {
		Name string `json:"username"`
		Pass string `json:"password"`
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		tool.Log("HandleLogin", "Sys", "localhost", "Controller Error")
	}
	decodedPassword, err := base64.StdEncoding.DecodeString(data.Pass)
	if err != nil {
		tool.Log("HandleLogin", "Sys", "localhost", "Controller Error")
	}
	success, uID, sessionID, name, class := tool.Validate(data.Name, string(decodedPassword))
	var result struct {
		Status    string `json:"status"`
		UID       string `json:"uID"`
		SessionID string `json:"sessionID"`
		Name      string `json:"name"`
	}
	if success {
		if !tool.HaveUser(uID) {
			succ := tool.NewUser(uID, name, class)
			if !succ {
				tool.Log("NewUser", "Sys", "localhost", "Database Error")
				result.Status = "fail"
				result.UID = ""
				b, err = json.Marshal(result)
				if err != nil {
					tool.Log("HandleLogin", "Sys", "localhost", "Controller Error")
				}
				w.Write(b)
				return
			}
		}

		var tempUser user
		tempUser.SessionID = sessionID
		users[uID] = tempUser
		result.Status, result.UID, result.SessionID, result.Name = "succ", uID, sessionID, name
		b, err = json.Marshal(result)
		if err != nil {
			tool.Log("HandleLogin", "Sys", "localhost", "Controller Error")
		}
		ip := r.Header.Get("X-Real-Ip")
		tool.Log("Login", result.UID, ip, "Succ")
		w.Write(b)
	} else {
		result.Status, result.UID, result.Name = "fail", "", ""
		b, err = json.Marshal(result)
		ip := r.Header.Get("X-Real-Ip")
		tool.Log("StudentValidate", "Unknown", ip, "Validate Fail")
		w.Write(b)
	}
}

func handleExit(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.Log("HandleExit", "Sys", "localhost", "Controller Error")
	}
	delete(users, string(b))
	ip := r.Header.Get("X-Real-Ip")
	tool.Log("HandleExit", string(b), ip, "Succ")
}
