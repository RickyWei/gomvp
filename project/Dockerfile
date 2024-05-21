# build
FROM golang:1.22-alpine3.19
COPY . /project
WORKDIR /project
RUN set -eux
RUN apk add make
RUN go env -w GOPROXY=https://goproxy.cn
RUN make build

# deploy
FROM alpine:latest
RUN set -eux
RUN apk add tzdata
# ENV TZ=Asia/Shanghai
COPY --from=0 /project/build /app
COPY --from=0 /project/conf.yaml /app/conf.yaml
COPY --from=0 /project/i18n/active.en.toml /app/active.en.toml
COPY --from=0 /project/i18n/active.zh.toml /app/active.zh.toml
WORKDIR /app
EXPOSE 80
CMD ["./app"]
