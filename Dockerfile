FROM golang:1.17-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./internal

FROM alpine:latest
WORKDIR /app
RUN apk add ffmpeg
ENV SDP_DIRECTORY "sdpcollection"
ENV DASH_SERVER "https://bifrost.inlive.app"
ENV MANIFEST_FILENAME "manifest.mpd"
ENV SDP_FILENAME "rtpforwarder.sdp"

ENV DB_HOST "localhost"
ENV DB_USER "admin"
ENV DB_PASS "admin"
ENV DB_NAME "livestream"

COPY --from=builder /app/main .
RUN mkdir logs && chmod -R 777 logs
RUN mkdir sdpcollection && chmod -R 777 sdpcollection

EXPOSE 8080
CMD ["./main"]
