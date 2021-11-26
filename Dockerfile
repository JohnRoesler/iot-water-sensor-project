FROM golang:1.17 AS build
COPY . /
WORKDIR /
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main

FROM scratch
COPY --from=build /main /usr/bin/main
EXPOSE 8888
ENTRYPOINT [ "/usr/bin/main" ]
