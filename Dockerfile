FROM golang:1.16-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN air

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/hello .

CMD ["./hello"]
