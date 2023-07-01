FROM golang:1.20
WORKDIR /app
COPY . .
RUN go build -o worker main.go

ENV WORKER=0
CMD ["sh", "-c", "/app/worker -WORKER=${WORKER}"]