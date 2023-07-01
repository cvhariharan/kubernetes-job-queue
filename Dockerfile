FROM golang:1.20
WORKDIR /app
COPY . .
RUN go build -o worker main.go

ENV WORKER=0
ENV EXIT=0

CMD ["sh", "-c", "/app/worker -WORKER=${WORKER} -EXIT=${EXIT}"]