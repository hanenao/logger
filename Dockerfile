FROM golang:1.17 as builder

WORKDIR /go-app

COPY . .

RUN git config --global url."https://$(cat private_repo_access_token.txt):x-oauth-basic@github.com/unipos".insteadOf "https://github.com/unipos"

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/start main.go

FROM alpine

COPY --from=builder /go/bin/start /start

CMD ["/start"]