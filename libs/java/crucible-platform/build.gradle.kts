// Thư viện dùng chung cho các dịch vụ Spring Boot của Crucible (M2, M3):
// sự kiện CloudEvents qua outbox và inbox (CR-DOC-03 mục 5), chuyển đổi token Keycloak (ADR-0004).
plugins {
    `java-library`
    id("io.spring.dependency-management") version "1.1.7"
}

group = "dev.crucible"
version = "0.1.0"

java {
    toolchain {
        // Mặc định Java 25 (LTS); -PjavaVersion=21 cho máy chưa có JDK 25.
        languageVersion = JavaLanguageVersion.of(providers.gradleProperty("javaVersion").getOrElse("25").toInt())
    }
}

repositories {
    mavenCentral()
}

dependencyManagement {
    imports {
        mavenBom("org.springframework.boot:spring-boot-dependencies:4.1.1")
    }
}

dependencies {
    api("org.springframework.boot:spring-boot-starter-amqp")
    api("org.springframework.boot:spring-boot-starter-jdbc")
    api("org.springframework.boot:spring-boot-starter-security-oauth2-resource-server")
    api("org.springframework.boot:spring-boot-starter-json")
    annotationProcessor("org.springframework.boot:spring-boot-configuration-processor")
}
