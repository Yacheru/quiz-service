FROM golang:alpine AS builder

WORKDIR /quiz-service

ADD ../go.mod .

COPY . .

RUN go build -o ./build/quiz.exe ./cmd/main.go

FROM alpine

WORKDIR /quiz-service

COPY ./configs/.env /quiz-service/configs/.env
COPY ./configs/config.json /quiz-service/configs/config.json
COPY ./migrations /quiz-service/migrations
COPY ./logs /quiz-service/logs
COPY ./web /quiz-service/web
COPY --from=builder /quiz-service/build /quiz-service

CMD ["/quiz-service/quiz.exe"]