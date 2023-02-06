FROM golang:1.19.4-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /restaurants

EXPOSE 9000

CMD [ "/restaurants" ]