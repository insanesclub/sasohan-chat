###################
### Build stage ###
###################

FROM    golang:1.15.8 AS builder
RUN     go get github.com/insanesclub/sasohan-chat
WORKDIR /go/src/github.com/insanesclub/sasohan-chat/
RUN     make build

###

FROM    fedora:33
WORKDIR /bin/
COPY    --from=builder /go/src/github.com/insanesclub/sasohan-chat/bin/chat .
EXPOSE  1323
CMD     ["./chat"]