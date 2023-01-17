# Start by building the application.
FROM golang:1.19 as build

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN swag init
RUN CGO_ENABLED=0 go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian11
COPY --from=build /go/bin/app /

EXPOSE 8000

CMD ["/app"]
