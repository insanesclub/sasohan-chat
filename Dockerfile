###################
### Build stage ###
###################

FROM golang:latest AS build

RUN go get github.com/insanesclub/sasohan-chat

WORKDIR /go/src/github.com/insanesclub/sasohan-chat
RUN make build

###

FROM alpine:3.13.1

COPY --from=build /go/src/github.com/insanesclub/sasohan-chat/bin/chat /bin

EXPOSE 1323

CMD ./bin/chat
