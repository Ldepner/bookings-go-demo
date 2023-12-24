FROM golang:1.19

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o bookings ./cmd/web/

EXPOSE 8080

ARG env
ENV envValue=$env

CMD ["sh", "-c", "./docker-run-${envValue}.sh"]
