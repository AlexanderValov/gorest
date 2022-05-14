FROM golang:1.17.0-alpine AS build

WORKDIR /go/src/app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/app ./cmd/main.go

FROM scratch
COPY --from=build /go/src/app/bin/app /bin/app
ENTRYPOINT ["/bin/app"]
EXPOSE 7008
