# Build stage
FROM golang:1.23.2-alpine3.20 AS build
# Update package lists and install ca-certificates
RUN apk update && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates
# RUN apk update && \ 
#     apk add ca-certificates
WORKDIR /app
COPY . .
RUN go build -o /server ./cmd/server

# Run stage  
# FROM scratch
FROM alpine:latest
WORKDIR /app
COPY --from=build /server /server
CMD ["/server"]
# CMD ["air", "-c", ".air.toml"]
