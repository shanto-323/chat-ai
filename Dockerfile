FROM golang:1.25.4-alpine3.22 AS build

WORKDIR /chat-ai
COPY . .

RUN go build -o /cmd/bin/app ./cmd


FROM alpine:3.22
WORKDIR /usr/bin
COPY --from=build /cmd/bin/app .
CMD [ "app" ]

