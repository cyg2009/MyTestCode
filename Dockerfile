#This is a debian flavor of our lambda base image.

FROM gcr.io/k8stest-167418/alpine-node:272212
RUN mkdir -p /var/runtime 
ADD runtime /var/runtime
ADD processrouter /var/runtime/processrouter
WORKDIR /var/runtime/func
ENV PATH="/var/runtime/rtsp/nodejs/bin:$PATH" \
    DEFAULT_SERVER_PORT="28903" \
    RUNTIME_ROOT="/var/runtime" \
    FUNCTION_REPOSITORY="http://10.21.119.117:5000/v2/serverless/"

EXPOSE 28903

CMD ["/var/runtime/processrouter"]

