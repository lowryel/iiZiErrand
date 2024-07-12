# build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git gcc libc-dev
WORKDIR /go/src/app
COPY . .
RUN go mod edit -module app
RUN go get -d -v ./...
RUN go install -v ./...

# final stage
FROM alpine:latest
LABEL Name=iiziErrand Version=0.0.1
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENV port=5432 dbname=iizidb user=eugene password=cartelo009 host=localhost
ENTRYPOINT ./app
EXPOSE 80

CMD [ "./app" ]
