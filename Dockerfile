FROM golang:alpine AS builder

RUN apk --update add git openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/bitbucket.org/kudoindonesia/microservice_user_management
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/vaksi/foreign_currency
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./


RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM gcr.io/distroless/base
COPY --from=builder /app /opt/
WORKDIR /opt
VOLUME /opt/configs
ENTRYPOINT ["/opt/app"]
