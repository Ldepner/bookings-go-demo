FROM golang:1.19

WORKDIR /usr/src/app

RUN go install -tags sqlite github.com/gobuffalo/pop/v6/soda@latest

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

EXPOSE 8080

COPY . .
RUN go build -o bookings ./cmd/web/

ARG env=development
RUN echo $env

RUN soda migrate -e ${env}

# This is overwritten on deployment for production
# ./bookings -dbname=bookings_ge86 -dbuser=lld -dbpass=mZS0Jcw6yAfriKCqRo7x17INg1Nvq4Pc -dbhost=dpg-cm45pqun7f5s73btffb0-a -cache=false
CMD ./bookings -dbname=bookings -dbuser=bytedance -dbhost=host.docker.internal -cache=false -production=false
