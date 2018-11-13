FROM golang:1.10 as builder

WORKDIR ${GOPATH}/src/github.com/dispatchframework/dispatch-events-cloudevents

COPY ["driver.go", "Gopkg.lock", "Gopkg.toml", "./"]

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /dispatch-events-cloudevents


FROM scratch

COPY --from=builder /dispatch-events-cloudevents /

ENTRYPOINT [ "/dispatch-events-cloudevents" ]
