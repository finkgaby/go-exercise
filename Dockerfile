FROM golang:1.18.2 AS builder

WORKDIR /app
RUN curl -o /etc/ssl/certs/dellca2018-bundle.crt -L --remote-name http://ecs-artifacts.cec.lab.emc.com/ovf/git/Dell2018/dellca2018-bundle.crt

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

# build the go application
RUN CGO_ENABLED=0 GOOS=linux go build -o Exercise .

FROM alpine:3.15.4

COPY --from=builder /app/Exercise .
COPY ./server/repositories/db/migrations/ ./server/repositories/db/migrations/

EXPOSE 8080

CMD ["./Exercise"]