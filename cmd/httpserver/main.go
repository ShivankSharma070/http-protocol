package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ShivankSharma070/http-protocol/internals/headers"
	"github.com/ShivankSharma070/http-protocol/internals/request"
	"github.com/ShivankSharma070/http-protocol/internals/response"
	"github.com/ShivankSharma070/http-protocol/internals/server"
)

const port = 42069

func toStr(data []byte) string{
	var out string = ""
	for _, b := range data{
		out+= fmt.Sprintf("%02x", b)
	}
	return out
}

func respond400() []byte {
	return []byte(` 
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
func respond500() []byte {
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
func respond200() []byte {
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
	server, err := server.Serve(port, func(w *response.Writer, r *request.Request) {
		h := response.GetDefaultHeaders(0)
		status := response.StatusOK
		body := respond200()
		if r.RequestLine.RequestTarget == "/yourproblem" {
			status = response.StatusBadRequest
			body = respond400()
		} else if r.RequestLine.RequestTarget == "/myproblem" {
			status = response.StatusInternalServerError
			body = respond500()
		} else if r.RequestLine.RequestTarget == "/video" {
			f, _ := os.ReadFile("assets/vim.mp4")
			w.WriteStatusLine(response.StatusOK)
			h.Replace("Content-type", "video/mp4")
			h.Replace("content-length",fmt.Sprintf("%d", len(f)))
			w.WriteHeaders(h)
			w.WriteBody(f)
		} else if strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin/") {
			target := r.RequestLine.RequestTarget
			res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])
			if err != nil {
				status = response.StatusInternalServerError
				body = respond500()
			} else{
				w.WriteStatusLine(response.StatusOK)
				h.Delete("content-length")
				h.Replace("content-type", "text/plain")
				h.Set("transfer-encoding", "chunked")
				h.Set("trailer", "X-Content-SHA256")
				h.Set("trailer", "X-Content-Length")
				w.WriteHeaders(h)

				fullBody := []byte{}
				for {
					data := make([]byte, 32)
					n, err:= res.Body.Read(data)
					if err != nil {
						break
					}

					fullBody = append(fullBody, data[:n]...)
					w.WriteBody([]byte(fmt.Sprintf("%X\r\n", n)))
					w.WriteBody(data[:n])
					w.WriteBody([]byte("\r\n"))
				}
				w.WriteBody([]byte("0\r\n"))
				trailer := headers.NewHeaders()
				out := sha256.Sum256(fullBody)
				trailer.Set("X-Content-SHA256", toStr(out[:]))
				trailer.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
				w.WriteHeaders(trailer)
				return 
			}
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
