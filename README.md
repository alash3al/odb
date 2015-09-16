ODB "Object DB"
===
> just a simple object-store, give it an object "image/document/buffer .." and it will return back with its position in the store
so you can store that position `key` on your own  
> you can create a simple `RESTful` layer on top of it then use it as a `RESTful Object Storage` .  
> you can build a simple cluster, just imagine

Why i built it ?
=====
> i want to re invent the wheel :3  
> i want to consume my free time  
> just for fun :3  
> to learn something new

Install
====
`go get github.com/alash3al/odb`

Usage
====

```go

package main

import(
	"fmt"
	"github.com/alash3al/odb"
	"bytes"
)

func main(){

  // open a [new] database
  // first param is the database file, second is the database file limit
  // it returns the `Database` pointer and an error "if found"
	db, e := odb.Open("./data.db", 1024 * 1024 * 2)

  // put a new value into the store
  // it accepts any object that has the `Read()` method
  // i.e: File.Read(), bytes.NewBuffer(), *http.Request.Body ... etc
  // it returns the position `k` and error `e` if found
	k, e := db.Put(bytes.NewBuffer([]byte("this is my value")))

  // fetch the data of the `k`
  // return "func()[]byte" and error "if found"
	next, err := db.Fetch(k)

  // recive the data ?
	data := ""
  
  // loop and fetch the data
	for {
		d := next()
		if d == nil {
			break
		} else {
			data += string(d)
		}
	}

	fmt.Println(data)
}

```
