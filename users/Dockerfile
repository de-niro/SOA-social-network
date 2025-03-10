FROM golang:1.24.1-alpine3.20
LABEL authors="de-niro"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./users -v .

RUN cp ./users /usr/local/bin/app

EXPOSE 8080

ENTRYPOINT ["app"]
#ENTRYPOINT ["top", "-b"]
