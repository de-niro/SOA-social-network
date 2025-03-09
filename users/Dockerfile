FROM golang:latest
LABEL authors="de-niro"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]
#ENTRYPOINT ["top", "-b"]