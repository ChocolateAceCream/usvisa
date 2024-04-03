FROM golang:alpine3.19 as build

WORKDIR /go/src/app
COPY ./main .
COPY ./config.yml .

FROM chromedp/headless-shell:latest
RUN apt update; apt install dumb-init -y; apt install libc6

ENTRYPOINT ["dumb-init", "--"]
COPY --from=build /go/src/app/main /tmp
COPY --from=build /go/src/app/config.yml /tmp
CMD ["/tmp/main"]

