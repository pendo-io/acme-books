package server

func Init(host string, port string) {
	r := NewRouter()
	r.RunOnAddr(host + ":" + port)
}
