run : build 
	@bin/tcpToHttp

build: 
	@go build  -o bin/tcpToHttp ./cmd/tcplistener

clean : 
	@rm -r bin

