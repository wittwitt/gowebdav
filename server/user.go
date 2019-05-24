package server

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/5dao/gdav/webdav"
)

// User web auth
type User struct {
	UID   string
	Pwd   string
	Root  string   // default .
	Hides []string // ref .gitingroe

	WebDav *webdav.Handler
}

//Init init user
func (user *User) Init() error {
	if !filepath.IsAbs(user.Root) {
		newRoot := filepath.Join(instancePath, user.Root)
		pathExist, err := PathExist(newRoot)
		if err != nil {
			log.Println(user.UID, newRoot, "PathExist err:", err)
			return err
		}
		if !pathExist {
			//todo
			//make dir
			log.Println(user.UID, newRoot, "Path not Exist ")
			return err
		}
		user.Root = newRoot
	}
	// /path/a/b not /path/a/b/
	user.Root = filepath.Join(user.Root, "")
	//todo
	//isDir

	log.Println(user.UID, "root", user.Root)

	return nil
}

//HandleLimits handle limit
// bool,is limit for path
func (user *User) HandleLimits(w http.ResponseWriter, r *http.Request) (isLimit bool, err error) {
	//
	log.Println(r.RequestURI, r.Method)

	// find prefix,,user power path
	// for path, power := range user.Powers {
	// 	if strings.Index(path, r.RequestURI) == 0 {

	// 	}
	// }

	// if !user.Powers[r.RequestURI].CheckLimit(r.Method) {
	// 	return false, nil
	// }

	//2 check Limits path exist and right

	return false, nil
}
