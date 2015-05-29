all: gchats gchatc

gchats: server.go
	go build -o gchats $^

gchatc: client.go
	go build -o gchatc $^

clean:
	rm gchats gchatc