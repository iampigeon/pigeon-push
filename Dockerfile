FROM golang:alpine

COPY /bin/pigeon-push /
WORKDIR /
EXPOSE 9030

CMD ["./pigeon-push"]
