FROM openzipkin/zipkin:2 as zipkin

FROM amazoncorretto:18-alpine

COPY --from=zipkin /usr/local/bin/ /usr/local/bin/

ARG USER=zipkin

WORKDIR /${USER}

RUN adduser -g '' -h ${PWD} -D ${USER}

COPY --from=zipkin --chown=${USER} /zipkin/ /zipkin/

USER ${USER}

RUN cd BOOT-INF/lib && wget https://github.com/aws/aws-msk-iam-auth/releases/download/v1.1.4/aws-msk-iam-auth-1.1.4-all.jar

ENTRYPOINT ["start-zipkin"]