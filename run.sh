#!/bin/bash

go build -o bookings cmd/web/*.go
./bookings -dbname=bookings -dbuser=bytedance -cache=false -production=false