FROM gcr.io/distroless/base

WORKDIR /app
COPY  ./dist/gobyoall-server-linux-amd64   ./gobyoall-server
CMD [ "/app/gobyoall-server" ]
