FROM golang:latest
RUN mkdir application
ADD /auth /application
RUN mkdir utils
ADD /utils /utils
WORKDIR /application/src
RUN go build -o app
CMD ["./app"]