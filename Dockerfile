#This is a debian flavor of our lambda base image.

#FROM gcr.io/k8stest-167418/alpine-node:272212
FROM mhart/alpine-node:8
RUN mkdir -p /var/runtime 
ADD runtime /var/runtime
ADD funcRouter /var/runtime/bin/funcRouter
WORKDIR /var/runtime/bin
ENV PATH="/var/runtime/bin:$PATH" \
    DEFAULT_SERVER_PORT="8903" \
    RUNTIME_ROOT="/var/runtime" \
    RUNTIME_LAMBDA="/var/runtime/bin/lambda-run"

EXPOSE 8903

CMD ["/var/runtime/bin/funcRouter"]

