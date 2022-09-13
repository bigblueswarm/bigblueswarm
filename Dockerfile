## Build
FROM golang:1.18-buster AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN make build

## Deploy
FROM gcr.io/distroless/base-debian11:nonroot

WORKDIR /

COPY --from=build /app/bin/b3lb-* /b3lb

EXPOSE 8090

USER nonroot:nonroot

ENTRYPOINT ["/b3lb"]

## docker run -it --mount type=bind,source="$(pwd)/config.yml",target=/config.yml,readonly -p 8090:8090 sledunois/b3lb:2.1.0 -config /config.yml