package server

//Config config
type Config struct {
	PassSalt string

	Listen string
	Prefix string

	StlCrt string //https
	StlKey string

	RootPath string

	Users map[string]*User
}
