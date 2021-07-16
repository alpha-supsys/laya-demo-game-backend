FROM alpine:3.10.2

RUN mkdir /app

COPY ./build/main /app/main

ENTRYPOINT [ "/app/main" ] 