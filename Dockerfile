FROM golang:1.21-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/cv-server .

FROM alpine:3.20

RUN addgroup -S cv && adduser -S -G cv cv
WORKDIR /app
COPY --from=build /out/cv-server /app/cv-server

# Read-only site and API. The container binds all interfaces so a published
# port works; publish it to host loopback only (127.0.0.1:8081:8081).
ENV CV_HOST=0.0.0.0 \
    CV_PORT=8081
EXPOSE 8081

USER cv:cv
CMD ["/app/cv-server"]
