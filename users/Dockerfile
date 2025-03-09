FROM golang:latest
LABEL authors="de-niro"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v ./...

RUN cp ./users /usr/local/bin/app

EXPOSE 8080

CMD ["app"]
#ENTRYPOINT ["top", "-b"]