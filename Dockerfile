
FROM golang:1.20-alpine AS builder
RUN apk update && apk add --no-cache make git gcc libc-dev linux-headers python3-dev
WORKDIR /nucleus
COPY . /nucleus
RUN make clean
RUN go mod tidy && go mod download
RUN make build
RUN python3 ./init.py -b /nucleus/build/nucleusd -c /root/.nucleus --overwrite 

FROM alpine:latest
WORKDIR /nucleus
COPY --from=builder /nucleus/build/nucleusd /usr/bin/nucleusd
COPY --from=builder /root/.nucleus /root/.nucleus
# COPY --from=builder /nucleus/nucleus_conf.yaml /root/.nucleus/config.yaml

CMD ["nucleusd", "start"]
