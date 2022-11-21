#[ BUILD ]######################################################################

# set up build image and dependencies
FROM golang:1.19.3-alpine AS build
WORKDIR /src
RUN apk add git

# download the code dependencies
COPY go.mod go.sum /src/
RUN go mod download

# copy the source code into the builder
COPY . .

# go install the binary as a static binary (including embedded copy of timezone
# database)
RUN CGO_ENABLED=0 go install -ldflags="-w -s" -tags timetzdata

#[ SCRATCH ]####################################################################

# set up runner image, based on no linux distribution
FROM scratch

# copy the statically-linked binary
COPY --from=build /go/bin/api /api

# install tls certificates:
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# configure run instructions
ENTRYPOINT ["/api"]
CMD ["serve"]
EXPOSE 80

