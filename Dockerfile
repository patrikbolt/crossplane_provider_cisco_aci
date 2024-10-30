# Build Stage
FROM golang:1.23 AS builder

WORKDIR /workspace

# Kopiere die Go-Moduldateien und lade die Abhängigkeiten
COPY go.mod go.sum ./
RUN go mod download

# Kopiere den Rest des Codes
COPY . ./

# Baue das Provider-Binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o provider cmd/provider/main.go

# Laufzeit-Image
FROM alpine:3.16

WORKDIR /

# Kopiere das gebaute Binary aus der Build-Stage
COPY --from=builder /workspace/provider .

# Kopiere package.yaml ins finale Image
COPY --from=builder /workspace/package.yaml .

# Setze das EntryPoint
ENTRYPOINT ["/provider"]

