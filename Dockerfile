FROM img.reg.3g:15000/ubuntu-base:v3
MAINTAINER xueying.zheng@yeepay.com
ADD conf /conf
#ADD Manifest /
ADD Dockerfile /
ADD main /sso

ENTRYPOINT ["/sso"]