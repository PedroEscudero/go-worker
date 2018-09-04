FROM golang:1.8

RUN mkdir /app
RUN echo $(pwd)
ADD . /app/
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]
