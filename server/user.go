package server

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"

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

//HandleLimits handle limit
// bool,is limit for path
func (user *User) HandleLimits(w http.ResponseWriter, r *http.Request) (isLimit bool, err error) {
	//
	//log.Println(r.RequestURI, r.Method)

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

//UserPwd UserPwd
func UserPwd(salt, uid, pwd string) string {
	var data []byte

	salt256 := sha256.Sum256([]byte(salt))
	data = append(data, salt256[:]...)

	uid256 := sha256.Sum256([]byte(salt + uid))
	data = append(data, uid256[:]...)

	pwd256 := sha256.Sum256([]byte(salt + uid + pwd))
	data = append(data, pwd256[:]...)

	sum := sha256.Sum256(data)

	return base64.StdEncoding.EncodeToString(sum[:])
}
