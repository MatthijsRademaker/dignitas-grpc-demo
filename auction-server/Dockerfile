FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /auction-server .

FROM gcr.io/distroless/static-debian12
COPY --from=build /auction-server /auction-server
EXPOSE 50051
ENTRYPOINT ["/auction-server"]
