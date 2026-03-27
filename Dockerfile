# syntax=docker/dockerfile:1

FROM golang:1.22 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=build /out/server /server

EXPOSE 8080
USER nonroot:nonroot

# Distroless images don't ship curl/wget; keep it simple here.
# If you want a real HTTP healthcheck, do it in docker-compose / Kubernetes.

ENTRYPOINT ["/server"]
