FROM alpine:latest AS build

RUN apk update
RUN apk upgrade
RUN apk add --update go=1.16.5-r0

WORKDIR /app
COPY . .

RUN go build -a -o out .


FROM alpine:latest

WORKDIR /app
COPY --from=build /app/out /app/out

ENTRYPOINT ["/app/out"]
