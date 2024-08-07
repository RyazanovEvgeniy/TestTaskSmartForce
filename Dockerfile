FROM golang:latest

COPY ./ ./
RUN go build -o main main.go
ENTRYPOINT [ "./main" ]
CMD ["-p", "./example/", "-o", "321", "-n", "3233"]