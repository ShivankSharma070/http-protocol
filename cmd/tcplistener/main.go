package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ShivankSharma070/http-protocol/internals/request"
)

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	out := make(chan string, 1)
//
// 	go func() {
// 		defer f.Close()
// 		defer close(out) // Closing the channel is important. Otherwise deadlock
//
// 		str := ""
// 		for {
// 			data := make([]byte, 8)
// 			n, err := f.Read(data)
// 			if err != nil {
// 				break
// 			}
//
// 			data = data[:n]
// 			if i := bytes.IndexByte(data, '\n'); i != -1 {
// 				str += string(data[:i])
// 				data = data[i+1:]
// 				out <- str
// 				str = ""
// 			}
//
// 			str += string(data)
// 		}
//
// 		if len(str) != 0 {
// 			out <- str
// 		}
//
// 	}()
// 	return out
// }

func main() {
	// Start a tcp connection
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("erorr", "error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
			break
		}

		re, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", "error", err)
			break
		}

		fmt.Println("Request line:")
		fmt.Println("- Method: ", re.RequestLine.Method)
		fmt.Println("- Target: ", re.RequestLine.RequestTarget)
		fmt.Println("- Version: ", re.RequestLine.HttpVersion)

	}

	listener.Close()
}
