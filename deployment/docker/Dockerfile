FROM golang:1.21-alpine as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server

FROM scratch as run
COPY --from=build /server .
CMD ["/server"]

