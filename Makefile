CLIENT=client
SERVER=server


all: $(CLIENT) $(SERVER)


clean:
	rm -f $(CLIENT) $(SERVER) 

$(CLIENT): chatmsg
	go build cmd/client.go
$(SERVER): chatmsg
	go build cmd/server.go