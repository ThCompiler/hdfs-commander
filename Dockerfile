FROM --platform=$BUILDPLATFORM golang:alpine

WORKDIR /hdfscmdr
COPY . .

RUN go build ./

CMD ["./hdfs-commander"]