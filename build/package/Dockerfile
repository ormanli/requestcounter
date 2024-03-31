FROM golang:1.23 as builder
RUN mkdir requestcounter
WORKDIR requestcounter

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -o requestcounter ./cmd/requestcounter/main.go

FROM scratch

COPY --from=builder /go/requestcounter/requestcounter .
CMD ["./requestcounter"]
