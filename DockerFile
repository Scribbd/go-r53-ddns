# BUILD STEP
FROM golang:1-alpine as go_build
WORKDIR /app
COPY ./golang ./
RUN go build ddns.go

# IMAGE ASSEMBLY
FROM alpine:latest
WORKDIR /app
COPY --from=go_build /app/ddns /app/ddns

ENTRYPOINT [ "/app/ddns" ]