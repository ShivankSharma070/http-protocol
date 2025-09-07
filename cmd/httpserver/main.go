package main

import (
	"log"
	"os"
	"os/signal"
	"fmt"
	"syscall"

	"github.com/ShivankSharma070/http-protocol/internals/request"
	"github.com/ShivankSharma070/http-protocol/internals/response"
	"github.com/ShivankSharma070/http-protocol/internals/server"
)

const port = 42069

func respond400() []byte{
	return []byte( ` 
	<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}
func respond500() []byte{
	return []byte(`
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}
func respond200() []byte{
	return []byte(`
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}

func main() {
	server, err := server.Serve(port, func(w *response.Writer, r *request.Request)  {
		h := response.GetDefaultHeaders(0)
		status := response.StatusOK
		body :=  respond200()
		if r.RequestLine.RequestTarget == "/yourproblem" {
			status = response.StatusBadRequest
			body = respond400()
		} else if r.RequestLine.RequestTarget == "/myproblem" {
			status = response.StatusInternalServerError
			body = respond500()
		}

		h.Replace("content-length", fmt.Sprintf("%d", len(body)))
		h.Replace("content-type", "text/html")
		w.WriteStatusLine(status)
		w.WriteHeaders(h)
		w.WriteBody(body)
	}) 

	if err != nil {
		log.Fatalf("Error Starting server: %v", err)
	}

	defer server.Close()

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server Gracefully stopped")
}
