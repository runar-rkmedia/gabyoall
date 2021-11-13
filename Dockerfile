FROM gcr.io/distroless/base

COPY  ./dist/gobyoall-server-linux-amd64 /gobyoall-server
CMD [ "/gobyoall-server" ]
