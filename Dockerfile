FROM golang:1.18-alpine AS builder
RUN apk update && apk add --no-cache make git gcc libc-dev linux-headers python3-dev
WORKDIR /nucleus
COPY . /nucleus
RUN rm -rf .git && git init
RUN go mod tidy && go mod download
RUN make clean build
RUN python3 ./init.py -b /nucleus/build/nucleusd -c /root/.nucleus --overwrite 
COPY addrbook.json /root/.nucleus/config/addrbook.json

FROM alpine:latest
WORKDIR /nucleus
COPY --from=builder /nucleus/build/nucleusd /usr/bin/nucleusd
COPY --from=builder /root/.nucleus /root/.nucleus
# COPY --from=builder /nucleus/nucleus_conf.yaml /root/.nucleus/config.yaml

CMD ["nucleusd", "start"]
