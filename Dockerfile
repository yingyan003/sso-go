FROM img.reg.3g:15000/ubuntu-base:v3
MAINTAINER xueying.zheng@yeepay.com
ADD Manifest /
ADD Dockerfile /
ADD src/main /sso

ENTRYPOINT ["/sso"]