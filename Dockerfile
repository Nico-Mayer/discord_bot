FROM golang:1.21.6-alpine

RUN apk add --no-cache ffmpeg --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community \
python3 \
py3-pip \
&& pip3 install --upgrade pip

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o out

ENTRYPOINT ["/app/out"]