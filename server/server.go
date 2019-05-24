package server

import (
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	"github.com/5dao/gdav/webdav"
)

var instancePath string
var sessionStore sessions.Store

func init() {
	var err error
	instancePath, err = GetInstancePath()
	if err != nil {
		log.Panicln("GetInstancePath err", err)
	}

	sessionStore = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
}

//NewServer make new webdav server
func NewServer(cfg *Config) (svr *Server, err error) {
	svr = &Server{
		Cfg:         cfg,
		Users:       make(map[string]*User),
		DavHandlers: make(map[string]*webdav.Handler),
	}

	//init user
	for _, user := range cfg.Users {
		if user.Init() != nil {
			return
		}
		svr.Users[user.UID] = user
		if _, ok := svr.DavHandlers[user.Root]; !ok {
			svr.DavHandlers[user.Root] = &webdav.Handler{
				Prefix:     cfg.Prefix + "/",
				FileSystem: webdav.Dir(user.Root),
				LockSystem: webdav.NewMemLS(),
				Hides:      user.Hides,
			}
		}
		user.WebDav = svr.DavHandlers[user.Root]

		//user.WebDav.MakeHides(user.Hides)
	}

	return
}

//Server dav server
type Server struct {
	Cfg *Config

	Users       map[string]*User
	DavHandlers map[string]*webdav.Handler
}

//Start start
func (svr *Server) Start() {
	go svr.Run()
}

//Run go run
func (svr *Server) Run() {
	defer func() {
		if rev := recover(); rev != nil {
			log.Println("Server run recover:", rev)
		}
		go svr.Run()
	}()

	mux := http.NewServeMux()
	mux.HandleFunc(svr.Cfg.Prefix+"/", svr.HandleFunc)

	err := http.ListenAndServe(svr.Cfg.Listen, mux)
	if err != nil {
		log.Println("runHTTP ListenAndServe err:", err)
	}
}

//HandleFunc handle req
//https://www.x.com/Prefix/path/a/b/c
func (svr *Server) HandleFunc(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)

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

	log.Println("......")
	user.WebDav.ServeHTTP(w, r)

	//webDavHandler.ServeHTTP(w, r)
}

//BasicAuth auth uid login
func (svr *Server) BasicAuth(w http.ResponseWriter, r *http.Request) *User {

	session, _ := sessionStore.Get(r, "gdav")
	if _, ok := session.Values["uid"]; ok {
		loginUID := session.Values["uid"].(string)
		return svr.Users[loginUID]
	}

	log.Println("not loin")

	uid, pwd, baseAuthOk := r.BasicAuth()
	if !baseAuthOk {
		w.Header().Set("WWW-Authenticate", `Basic realm=""`)
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}
	log.Println("BasicAuth ok", uid, pwd)

	var loginUser *User

	checkUserOk := false
	usersLen := len(svr.Cfg.Users)
	for i := 0; i < usersLen; i++ {
		//todo
		//pwd = mad5(pwd+google code)
		if svr.Cfg.Users[i].UID == uid && svr.Cfg.Users[i].Pwd == pwd {
			checkUserOk = true
			loginUser = svr.Cfg.Users[i]
		}
	}
	if !checkUserOk {
		w.Header().Set("charset", "UTF-8")
		//firefox realm chinese unintelligible text
		w.Header().Set("WWW-Authenticate", `Basic realm="UID/PWD Error!"`)
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	log.Println("checkUserOk ok", uid, pwd)

	return loginUser
}
