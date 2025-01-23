FROM golang:1.23-alpine AS build

#RUN useradd -u 1001 nonroot

WORKDIR /app 

COPY go.mod ./

RUN go mod download

COPY . ./

RUN go build -ldflags="-s -w" -o hello-app ./

FROM alpine

WORKDIR /app

COPY --from=build /app/hello-app /app/hello-app

# Use the non-root user created in the build stage TODO
#USER nonroot

CMD ["/app/hello-app"]

