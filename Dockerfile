# Verwende Go 1.23 als Basisversion
FROM golang:1.23 AS builder

# Setze das Arbeitsverzeichnis
WORKDIR /workspace

# Kopiere die Go-Module und lade die Abhängigkeiten
COPY go.mod go.sum ./
RUN go mod download

# Kopiere den restlichen Code
COPY . .

# Baue das Provider-Binary mit dem korrekten Pfad
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o provider ./cmd/provider/main.go

# Verwende ein leichtes Basis-Image für die finale Version
FROM alpine:3.16
WORKDIR /

# Kopiere das Binary aus dem Build-Container
COPY --from=builder /workspace/provider .

# Setze den Startpunkt des Containers
ENTRYPOINT ["/provider"]

