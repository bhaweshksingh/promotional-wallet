FROM golang:alpine as builder
WORKDIR /account
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o account cmd/*.go

FROM scratch
COPY --from=builder /account/account .
ENTRYPOINT ["./account","http-serve"]