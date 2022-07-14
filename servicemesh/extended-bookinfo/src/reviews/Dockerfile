FROM java:8

ARG service_version
ARG enable_ratings
ARG star_color

ENV SERVICE_VERSION ${service_version:-v1}
ENV ENABLE_RATINGS ${enable_ratings:-false}
ENV STAR_COLOR ${star_color:-black}

COPY target/reviews-0.1.0-SNAPSHOT.jar /app/reviews.jar

CMD["java", "-jar", "/app/reviews.jar"]