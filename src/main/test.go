package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"main/session"
	_ "main/session/providers/memory"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type Friend struct {
	Fname string
}

type Person struct {
	UserName string
	Emails   []string
	Friends  []*Friend
}

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
}

func EmailDealWith(args ...interface{}) string {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}
	// find the @ symbol
	substrs := strings.Split(s, "@")
	if len(substrs) != 2 {
		return s
	}
	// replace the @ by " at "
	return (substrs[0] + " at " + substrs[1])
}

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("Path", r.URL.Path)
	fmt.Println("Scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("Key:", k)
		fmt.Println("Val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello Czjge!")
}

func login(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	r.ParseForm()
	fmt.Println(`Method:`, r.Method)
	if r.Method == "GET" {
		// crutime := time.Now().Unix()
		// h := md5.New()
		// io.WriteString(h, strconv.FormatInt(crutime, 10))
		// token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("login.gtpl")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, sess.Get("username"))
	} else {
		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/", 302)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`Method:`, r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		f1 := Friend{Fname: "minux.ma"}
		f2 := Friend{Fname: "xushiwei"}
		p := Person{UserName: "Astaxie",
			Emails:  []string{"astaxie@beego.me", "astaxie@gmail.com"},
			Friends: []*Friend{&f1, &f2}}

		//t, _ := template.ParseFiles("upload.gtpl")
		//t.Funcs(template.FuncMap{"bar": EmailDealWith}).Execute(w, map[string]interface{}{"token": token, "user": p})

		t, _ := template.ParseFiles("tmpl/header.tmpl", "tmpl/content.tmpl", "tmpl/footer.tmpl")
		t.ExecuteTemplate(w, "header", map[string]interface{}{"token": token, "user": p})
		t.ExecuteTemplate(w, "content", map[string]interface{}{"token": token, "user": p})
		t.ExecuteTemplate(w, "footer", map[string]interface{}{"token": token, "user": p})
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func count(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Query()["a"])

	sess := globalSessions.SessionStart(w, r)
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	t, _ := template.ParseFiles("count.gtpl")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("countnum"))
}

func message(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("message.html")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, nil)
}

func Echo(ws *websocket.Conn) {
	var err error

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received back from client: " + reply)

		msg := "Received: " + reply
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

func main() {

	// arith := new(Arith)
	// rpc.Register(arith)
	// rpc.HandleHTTP()

	// err := http.ListenAndServe(":1234", nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// http.Handle("/", websocket.Handler(Echo))
	// if err := http.ListenAndServe(":1234", nil); err != nil {
	// 	log.Fatal("ListenAndServe:", err)
	// }

	http.HandleFunc("/", sayHelloName)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/count", count)
	http.HandleFunc("/message", message)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
