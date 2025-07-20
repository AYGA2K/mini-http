package server

type Handler interface {
	Serve(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) Serve(w ResponseWriter, r *Request) {
	f(w, r)
}
