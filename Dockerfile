FROM crystallang/crystal:1.1.1-alpine-build as builder

WORKDIR /usr/src/app

COPY . .

RUN shards install --production

RUN shards build --release --no-debug -s

FROM crystallang/crystal:1.1.1-alpine as runner

RUN apk add dumb-init

WORKDIR /usr/src/app

COPY --from=builder ["/usr/src/app/bin", "."]

ENV HOST="0.0.0.0"

ENV PORT="3000"

EXPOSE 3000

CMD ["dumb-init", "./main"]