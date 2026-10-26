# Build trong thư mục gốc kho mã: docker build -f services/<module>/Dockerfile .
# (cần contracts/ và libs/ để sinh mã và biên dịch)
FROM eclipse-temurin:25-jdk AS build
WORKDIR /src
COPY contracts contracts
COPY libs/java libs/java
COPY services/learning services/learning
WORKDIR /src/services/learning
RUN ./gradlew --no-daemon bootJar -x test

FROM eclipse-temurin:25-jre
RUN useradd --system --uid 10001 app
USER 10001
WORKDIR /app
COPY --from=build /src/services/learning/build/libs/*-SNAPSHOT.jar app.jar
EXPOSE 8080
ENTRYPOINT ["java", "-XX:MaxRAMPercentage=75", "-jar", "/app/app.jar"]
