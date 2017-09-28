FROM gcr.io/k8stest-167418/alpine-node:272212

ADD processrouter /bin/processrouter

CMD [ "/bin/processrouter" ]




