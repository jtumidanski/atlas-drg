FROM maven:3.6.3-openjdk-14-slim AS build

COPY settings.xml /usr/share/maven/conf/

COPY pom.xml pom.xml
COPY drg-api/pom.xml drg-api/pom.xml
COPY drg-model/pom.xml drg-model/pom.xml
COPY drg-base/pom.xml drg-base/pom.xml

RUN mvn dependency:go-offline package -B

COPY drg-api/src drg-api/src
COPY drg-model/src drg-model/src
COPY drg-base/src drg-base/src

RUN mvn install

FROM openjdk:14-ea-jdk-alpine
USER root

RUN mkdir service

COPY --from=build /drg-base/target/ /service/
COPY config.yaml /service/

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.5.0/wait /wait

RUN chmod +x /wait

ENV JAVA_TOOL_OPTIONS -agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005

EXPOSE 5005

CMD /wait && java --enable-preview -jar /service/drg-base-1.0-SNAPSHOT.jar -Xdebug