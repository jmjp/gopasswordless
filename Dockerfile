#Start by building the application.
FROM golang:1.21.5-bullseye as build

WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app main.go


# Now copy it into our base image.
FROM scratch
# FROM scratch
COPY --from=build /go/bin/app /

CMD ["/app"]

