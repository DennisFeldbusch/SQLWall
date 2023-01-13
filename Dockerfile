FROM golang:1.19

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /usr/local/bin/app

# Incoming
EXPOSE 8080 
# Outgoing
EXPOSE 80
CMD ["/usr/local/bin/app"]