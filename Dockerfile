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

ENV DB_HOST "34.101.236.119"
ENV DB_USER "postgres"
ENV DB_PASS "admin"
ENV DB_NAME "kZ3yr37JLY3Nr3JUbuEg"

COPY --from=builder /app/main .
RUN mkdir logs && chmod -R 777 logs
RUN mkdir sdpcollection && chmod -R 777 sdpcollection


EXPOSE 8080
CMD ["./main"]
