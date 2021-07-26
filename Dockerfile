FROM golang:1.16-alpine AS build
WORKDIR /app
COPY . .

# for go binary to run on alpine image 
RUN CGO_ENABLED=0 go build -o /shorturl

FROM alpine:latest AS deploy
WORKDIR /
COPY --from=build /shorturl /
ENV GIN_MODE=release \
	PORT=8080

CMD ["/shorturl"]
