# micro-go-chat
Sample project to play with go, websockets, CI and Docker

Usual way to run a go thing:

    go get github.com/aveius/micro-go-chat
    go build github.com/aveius/micro-go-chat
    go run github.com/aveius/micro-go-chat

After that, open a couple browser windows on http://127.0.0.1:8080/, and you're off!

A local postgre instance can be used to add some persistency to this. To leverage it, just configure `PG_CONN_URL` – see the "Connection String Parameters" in [`pq`'s documentation](https://godoc.org/github.com/lib/pq). For instance:

    postgres://USER:PASSWORD@127.0.0.1:5432/?sslmode=disable

By default, access is limited to 127.0.0.1 (**/!\ No X-Forwarded checks**); if deploying to a remote server, put your connecting IP in the ENVVAR `WHITELISTED_IP`.


## Docker packaging and run

After checking out the project:

    make build
    sudo make docker

And the usual to spin it up:

    docker run -p 8080:8080 -it -e "WHITELISTED_IP=172.17.0.1" -e "PG_CONN_URL=postgres://..." micro-go-chat



## Status
- [x] websockets
- [x] Persistency with postgre
- [x] Docker deployment
- [ ] CI/CD
