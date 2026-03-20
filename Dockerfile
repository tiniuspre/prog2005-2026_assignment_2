FROM ubuntu:latest
LABEL authors="stankel"

ENTRYPOINT ["top", "-b"]