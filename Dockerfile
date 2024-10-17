# Build stage
FROM golang:1.23.2-alpine3.20 AS build
WORKDIR /app
COPY . .
RUN go build -o /server ./cmd/server

# Run stage  
FROM scratch
WORKDIR /app
COPY --from=build /server /server
CMD ["/server"]
# CMD ["air", "-c", ".air.toml"]
