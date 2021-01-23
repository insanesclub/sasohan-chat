###################
### Build stage ###
###################

FROM golang:1.15.7 AS build

MAINTAINER msh0117@kookmin.ac.kr

WORKDIR /go/src/github.com/insanesclub
RUN go get github.com/insanesclub/sasohan-chat

WORKDIR /go/src/github.com/insanesclub/sasohan-chat
RUN make build

###

FROM fedora:33

COPY --from=build /go/src/github.com/insanesclub/sasohan-chat/bin/chat /bin

EXPOSE 1323

CMD /bin/chat
