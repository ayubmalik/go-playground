FROM golang:1.23-alpine AS build

#RUN useradd -u 1001 nonroot
#RUN addgroup -S hello-app && adduser -S hello-app -G hello-app

#USER hello-app

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . ./

RUN go test -timeout 1m -v ./... > test.log

RUN ls

RUN go build -ldflags="-s -w" -o /app/hello-app ./

FROM alpine

WORKDIR /app

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /app/hello-app /app/hello-app

# Use the non-root user created in the build stage TODO
#USER nonroot
#USER hello-app

CMD ["/app/hello-app"]

