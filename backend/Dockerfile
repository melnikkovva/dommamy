FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server
RUN mkdir -p /app/uploads

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/server /server
COPY --from=builder --chown=65532:65532 /app/uploads /app/uploads
ENV FILES_DIR=/data/uploads
EXPOSE 8080
ENTRYPOINT ["/server"]