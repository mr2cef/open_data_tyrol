FROM golang:1.16-buster

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /collect

EXPOSE 8080

CMD [ "/collect" ]