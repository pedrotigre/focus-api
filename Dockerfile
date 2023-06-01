FROM golang:1.20.4-alpine3.18 AS build

WORKDIR /app

COPY . /app

ARG RAILWAY_ENVIRONMENT

RUN go mod download && go mod tidy

RUN go build -o focus-api .

# Final stage
FROM alpine:3.14

WORKDIR /app

COPY --from=build /app/focus-api .

ENV RAILWAY_ENVIRONMENT=$RAILWAY_ENVIRONMENT

EXPOSE 8080

CMD [ "./focus-api" ]