package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"sync"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"

	"github.com/5dao/gdav/webdav"
)

//Server gdav server
type Server struct {
	Cfg *Config

	DavHandlers map[string]*webdav.Handler

	Users *sync.Map //[uid]User

	Router *httprouter.Router

	Sesstions *sessions.FilesystemStore
}

//NewServer make gdav server
func NewServer(cfg *Config) (svr *Server, err error) {
	if cfg.RootPath == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		cfg.RootPath = filepath.Join(usr.HomeDir, ".gdav")
	}
	cfg.RootPath, err = filepath.Abs(cfg.RootPath)
	if err != nil {
		return nil, err
	}

	svr = &Server{
		Cfg: cfg,

		DavHandlers: make(map[string]*webdav.Handler),
		Users:       new(sync.Map), //,

		Router: httprouter.New(),

		Sesstions: sessions.NewFilesystemStore(cfg.RootPath+"/.sesstions", securecookie.GenerateRandomKey(32)),
	}

	//make webdav
	for uid, user := range cfg.Users {
		user.UID = uid

		// user root
		if !filepath.IsAbs(user.Root) {
			user.Root = filepath.Join(cfg.RootPath, user.Root)
		}

		rootInfo, statErr := os.Stat(user.Root)
		if os.IsNotExist(statErr) {
			log.Printf("path not exist: %s", user.Root)
			err = os.MkdirAll(user.Root, 0770)
			if err != nil {
				return nil, err
			}
			rootInfo, statErr = os.Stat(user.Root)
		}
		if statErr != nil {
			return nil, fmt.Errorf("%s err: %v", user.Root, statErr)
		}
		if !rootInfo.IsDir() {
			return nil, fmt.Errorf("%s is file not dir", user.Root)
		}

		if _, ok := svr.DavHandlers[user.Root]; !ok {
			userHides := []string{}
			for _, hidePath := range user.Hides {
				userHides = append(userHides, "/"+cfg.Prefix+"/"+hidePath)
			}

			svr.DavHandlers[user.Root] = &webdav.Handler{
				Prefix:     "/" + cfg.Prefix,
				FileSystem: webdav.Dir(user.Root),
				LockSystem: webdav.NewMemLS(),
				Hides:      userHides, // user.Hides,
			}
		}
		user.WebDav = svr.DavHandlers[user.Root]

		svr.Users.Store(user.UID, user)
	}

	svr.RegRouter()

	return
}

//Start start
func (svr *Server) Start() {
	go svr.Run()
}

//Run go run
func (svr *Server) Run() {
	// defer func() {
	// 	if rev := recover(); rev != nil {
	// 		log.Println("Server run recover:", rev)
	// 	}
	// 	go svr.Run()
	// }()

	if svr.Cfg.StlCrt != "" {
		err := http.ListenAndServeTLS(svr.Cfg.Listen, svr.Cfg.StlCrt, svr.Cfg.StlKey, svr.Router)
		if err != nil {
			log.Println("runHTTPs ListenAndServe err:", err)
		}
	} else {
		err := http.ListenAndServe(svr.Cfg.Listen, svr.Router)
		if err != nil {
			log.Println("runHTTP ListenAndServe err:", err)
		}
	}
}

// RegRouter RegRouter
func (svr *Server) RegRouter() {
	webDavPrefix := "/" + svr.Cfg.Prefix + "/*filepath"

	// WebDAV
	for _, method := range WebDAVMethods {
		svr.Router.HandlerFunc(method, webDavPrefix, svr.HandleFunc)
	}

	svr.Router.NotFound = http.HandlerFunc(NotFound)

}

//HandleFunc handle req
func (svr *Server) HandleFunc(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.RequestURI)

	// auth user
	user := svr.BasicAuth(w, r)
	if user == nil {
		return
	}

	isLimit, err := user.HandleLimits(w, r)
	if err != nil {
		log.Println("HandleFunc:", err)
	}
	if isLimit {
		return
	}

	user.WebDav.ServeHTTP(w, r)
}

//BasicAuth auth uid login
func (svr *Server) BasicAuth(w http.ResponseWriter, r *http.Request) *User {
	gdavSession, err := svr.Sesstions.Get(r, "gdav")
	if err != nil {
		log.Println(err)
		return nil
	}

	if userObj, ok := gdavSession.Values["uid"]; ok {
		return userObj.(*User)
	}

	//
	uid, pwd, baseAuthOk := r.BasicAuth()
	if !baseAuthOk {
		w.Header().Set("WWW-Authenticate", `Basic realm=""`)
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	var loginUser *User
	if userObj, ok := svr.Users.Load(uid); ok {
		loginUser = userObj.(*User)
	}

	w.Header().Set("charset", "UTF-8")
	//firefox realm chinese unintelligible text

	if loginUser == nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="UID/PWD Error! code=101"`)
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	if UserPwd(svr.Cfg.PassSalt, uid, pwd) != loginUser.Pwd {
		w.Header().Set("WWW-Authenticate", `Basic realm="UID/PWD Error! code=102"`)
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	gdavSession.Options.MaxAge = 0

	gdavSession.Values["uid"] = loginUser
	svr.Sesstions.Save(r, w, gdavSession)

	return loginUser
}

// NotFound 404 action
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNotFound) // StatusNotFound = 404
	w.Write([]byte("My own Not Found handler." + r.URL.String()))
	w.Write([]byte(" The page you requested could not be found."))
}
