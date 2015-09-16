// no comment :3
package odb

// fetch the required packages
import(
	"os"
	"io"
	"sync"
	"errors"
	"fmt"
)

// database is the main container
type Database struct {
	sync.RWMutex
	file		*os.File
	size 		int64
	limit 		int64
	full 		bool
}

// create a new instance of the database
func Open(filename string, limit int64) (_ *Database, err error) {

	// initialize a new Database
	// and set its initial valus
	this := new(Database)
	this.file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	finfo, _ := this.file.Stat()
	this.size = finfo.Size()
	this.limit = limit

	// check for any error
	// and some other configs
	if err != nil {
		return this, err
	}

	// everything is ok
	return this, nil

}

// close the database file
func (this *Database) Close(){
	this.file.Close()
	this = nil
}

// is the database full ?
func (this *Database) Full() bool {
	return (this.limit - 4096) < this.size
}

// return the current size of the database
func (this *Database) Size() int64 {
	return this.size
}

// put a new object into the store and get its position
func (this *Database) Put(src io.Reader) (string, error) {

	// the database is full ?
	if this.Full() {
		return "", errors.New("the database is full")
	}

	// lock the database
	this.Lock()
	defer this.Unlock()

	// now copy the data from the object to our database
	// then ensure that there is no error
	size, err := io.Copy(this.file, src)
	if err != nil {
		return "", err
	}

	// the data offset
	offset := this.size

	// increment the database size
	this.size += size

	// finalize, [offset:size]
	return fmt.Sprintf("%d:%d", offset, size), nil

}

// fetch an object from the store using its position
func (this *Database) Fetch(pos string, cb func([]byte)) error {

	// lock the database "for-reading"
	this.RLock()
	defer this.RUnlock()

	// prepare the required vars for parsing
	var offset, size int64

	// parse the pos
	fmt.Sscanf(pos, "%d:%d", &offset, &size)

	// is valid offset ?
	if this.size <= offset {
		return errors.New("Invalid offset")
	}

	// is valid size ?
	if (size < 1) && (size >= this.size) {
		return errors.New("Invalid offset")
	}

	// declare tha variables that used each iteration
	next := offset
	reminder := size
	chunk := int64(1024 * 32)

	// start copying the data
	for {

		if (reminder < chunk) {
			chunk = reminder
		}

		buf := make([]byte, chunk)

		r, e := this.file.ReadAt(buf, next)

		if (e != nil) || (r == 0) {
			cb(nil)
			break
		}

		cb(buf[:r])

		reminder -= int64(r)
		next += int64(r)
	}

	// finalize
	return nil

}
