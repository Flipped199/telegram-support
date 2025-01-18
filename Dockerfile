FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

RUN apk add --no-cache --update build-base wget

RUN wget -O sqlite.tar.gz https://www.sqlite.org/src/tarball/sqlite.tar.gz?r=release && \
    mkdir -p /sqlite-src && \
    tar -xzf sqlite.tar.gz -C /sqlite-src --strip-components=1 && \
    cd /sqlite-src && \
    ./configure --prefix=/usr && make && make install

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bot .

FROM alpine:latest

RUN apk --no-cache add tzdata libc6-compat libgcc libstdc++

WORKDIR /app

COPY --from=builder /app/bot .

RUN mkdir -p /app/conf

ENTRYPOINT ["./bot","-c","./conf/config.toml"]