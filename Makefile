run : build 
	@bin/tcpToHttp

build: 
	@go build -o bin/tcpToHttp

clean : 
	@rm -r bin

