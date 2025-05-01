FROM golang:1.24.2-alpine3.20 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server

FROM scratch as run
COPY --from=build /server .
CMD ["/server"]

