package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type expression struct {
	Exp string  `json:"exp"`
	Res float64 `json:"res"`
}

type Stack []string

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(exp string) {
	*s = append(*s, exp)
}

func (s *Stack) Pop() string {

	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element

}

func (s *Stack) Peek() string {

	index := len(*s) - 1
	element := (*s)[index]
	return element

}

func check(op1, op2 string) bool {

	if (op1 == "*" || op1 == "/") && (op2 == "+" || op2 == "-") {
		return false
	} else {
		return true
	}
}

func cal(op string, bnum string, anum string) (float64, bool) {

	a, _ := strconv.ParseFloat(anum, 64)
	b, _ := strconv.ParseFloat(bnum, 64)

	switch op {
	case "+":
		return a + b, true
	case "-":
		return a - b, true
	case "*":
		return a * b, true
	case "/":
		if b == 0 {
			return -1, false
		}
		return a / b, true
	}
	return 0, false
}

func fetch(w http.ResponseWriter, r *http.Request) {

	var data expression
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	var values, ops Stack

	for i := 0; i < len(data.Exp); i++ {

		val := data.Exp[i]
		r := string(val)

		if (val >= 48 && val <= 57) || r == "." {

			sbt := ""
			for i < len(data.Exp) {

				val = data.Exp[i]
				r = string(val)

				if (val >= 48 && val <= 57) || r == "." {
					sbt = sbt + string(val)
				} else {
					break
				}
				i++
			}
			i--
			values.Push(sbt)

		} else if r == "+" || r == "-" || r == "*" || r == "/" {

			for !ops.IsEmpty() && check(r, ops.Peek()) {
				num, err := cal(ops.Pop(), values.Pop(), values.Pop())
				if !err {
					response, _ := json.Marshal(struct {
						Res string `json:"res"`
					}{"divide by zero"})
					w.Write(response)
					return
				}
				values.Push(fmt.Sprint(num))
			}
			ops.Push(r)
		}
	}
	for !ops.IsEmpty() {

		num, err := cal(ops.Pop(), values.Pop(), values.Pop())
		if !err {
			response, _ := json.Marshal(struct {
				Res string `json:"res"`
			}{"divide by zero"})
			w.Write(response)
			return
		}
		values.Push(fmt.Sprint(num))
	}

	data.Res, _ = strconv.ParseFloat(values.Pop(), 64)

	file_content, _ := json.MarshalIndent(data, "", "	")

	f, err := os.OpenFile("Information.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(file_content)
	if err != nil {
		f.Close()
		log.Fatal(err)
	}
	f.Close()

	w.Write(file_content)
}

func main() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/fetch", fetch)
	fmt.Println("listening 8080")

	http.ListenAndServe(":8080", nil)
}
