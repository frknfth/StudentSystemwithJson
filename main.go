package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Page is the form of txt file
type Page struct {
	Title string
	Body  []byte
}

// Student for json type
type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Uni  string `json:"uni"`
	ID   int    `json:"id"`
}

func (p Student) toString() string {
	return toJSON(p)
}
func toJSON(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(bytes)
}
func getPages() []Student {
	raw, err := ioutil.ReadFile("./students.txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var c []Student
	json.Unmarshal(raw, &c)
	return c
}

var lastline int

func (p *Page) saveToTxtFile() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}
func toPagefromTxtFile(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
func readTxtFile(title string) string {
	slice, err := ioutil.ReadFile(title + ".txt")
	if err != nil {
		return ""
	}
	return string(slice) + "\n"
}
func test(rw http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	checkErr(err)
	student := Student{}
	json.Unmarshal(body, &student)

	requestedID := student.ID

	pages := getPages()
	for _, p := range pages {
		if p.ID == requestedID {
			rw.Write([]byte("There is another student with same id"))
			return
		}
	}
	if lastline == 1 {
		p1 := &Page{Title: "students", Body: []byte(readTxtFile("students") + "\t" + string(body) + "\n]")}
		p1.saveToTxtFile()
		toPagefromTxtFile("students")
		lastline = lastline + 2
	} else {
		removeLines("students.txt", lastline, 1)
		p1 := &Page{Title: "students", Body: []byte(readTxtFile("students") + "\t,\n\t" + string(body) + "\n]")}
		p1.saveToTxtFile()
		toPagefromTxtFile("students")
		lastline = lastline + 3
	}
}

func main() {
	p1 := &Page{Title: "students", Body: []byte("[")}
	p1.saveToTxtFile()
	toPagefromTxtFile("students")
	lastline = lastline + 1

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/homePage.html")
	})

	http.HandleFunc("/homePagesjs.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/homePagesjs.js")
	})

	http.HandleFunc("/findPage", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/findPage.html")
	})

	http.HandleFunc("/findPagejs.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/findPagejs.js")
	})

	http.HandleFunc("/design.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/design.css")
	})

	http.HandleFunc("/save", test)

	http.HandleFunc("/get", get)

	log.Fatal(http.ListenAndServe(":1111", nil))

}

func get(rw http.ResponseWriter, req *http.Request) {
	requestedIDStr := req.Header.Get("Id")
	requestedID, _ := strconv.Atoi(requestedIDStr)
	var requestedStudent Student

	pages := getPages()
	for _, p := range pages {
		if p.ID == requestedID {
			requestedStudent = p
			break
		}
	}
	json.NewEncoder(rw).Encode(requestedStudent)
}

func removeLines(fn string, start, n int) (err error) {
	if start < 1 {
		return errors.New("invalid request.  line numbers start at 1")
	}
	if n < 0 {
		return errors.New("invalid request.  negative number to remove")
	}
	var f *os.File
	if f, err = os.OpenFile(fn, os.O_RDWR, 0); err != nil {
		return
	}
	defer func() {
		if cErr := f.Close(); err == nil {
			err = cErr
		}
	}()
	var b []byte
	if b, err = ioutil.ReadAll(f); err != nil {
		return
	}
	cut, ok := skip(b, start-1)
	if !ok {
		return fmt.Errorf("less than %d lines", start)
	}
	if n == 0 {
		return nil
	}
	tail, ok := skip(cut, n)
	if !ok {
		return fmt.Errorf("less than %d lines after line %d", n, start)
	}
	t := int64(len(b) - len(cut))
	if err = f.Truncate(t); err != nil {
		return
	}
	if len(tail) > 0 {
		_, err = f.WriteAt(tail, t)
	}
	return
}
func skip(b []byte, n int) ([]byte, bool) {
	for ; n > 0; n-- {
		if len(b) == 0 {
			return nil, false
		}
		x := bytes.IndexByte(b, '\n')
		if x < 0 {
			x = len(b)
		} else {
			x++
		}
		b = b[x:]
	}
	return b, true
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
