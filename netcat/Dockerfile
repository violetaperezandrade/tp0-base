FROM alpine:latest
RUN apk add --update netcat-openbsd
COPY netcat/netcat_script.sh  /
RUN chmod +x /netcat_script.sh
CMD nc -l -p 1234