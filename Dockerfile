FROM golang:1.15.6-alpine3.12 as base
EXPOSE 8080


FROM golang:1.15.6-alpine3.12 as build
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN go mod download
RUN go build -o server .


FROM base as final
RUN mkdir /app
RUN mkdir /app/storage
WORKDIR /app
COPY --from=build /build/server .
VOLUME storage
CMD ["/app/server"]
