FROM golang:1.19

COPY . /github.com/vaberof/TelegramBotUniversitySchedule/
WORKDIR /github.com/vaberof/TelegramBotUniversitySchedule/

RUN go mod download
RUN go build -o ./bin/app cmd/app/main.go

EXPOSE 80

CMD ["./bin/app"]
