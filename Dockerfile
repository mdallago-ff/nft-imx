##############
# Build Step #
##############
ARG BASE_IMAGE=golang:1.19.0-alpine3.16

FROM ${BASE_IMAGE} AS builder

ARG NAME=nft-imx
ARG SERVERS_DIR=.

RUN apk add --no-cache build-base git curl wget libressl-dev openssh

WORKDIR /app

# Copy everything
COPY ${NAME}/. .

# Download Go mods
RUN go mod download && go mod verify

# Build the project
RUN go build -v -o /${NAME} ./${SERVERS_DIR}

###############
# Deploy Step #
###############

FROM ${BASE_IMAGE} AS candidate

ARG NAME=nft-imx
ENV NAME=$NAME
ENV WORKER_NAME=$WORKER_NAME

# Install nice to haves
RUN apk add --no-cache openssl ncurses-libs libstdc++ libgcc curl libressl

WORKDIR /app

RUN chown nobody:nobody /app

USER nobody:nobody

# Copy everything over from the builder
COPY --from=builder --chown=nobody:nobody /$NAME .

# Expose 4000, set Entrypoint for app
EXPOSE 4000

CMD ["./nft-imx"]
