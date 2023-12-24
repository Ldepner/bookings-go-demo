FROM ubuntu:latest
LABEL authors="bytedance"

ENTRYPOINT ["top", "-b"]