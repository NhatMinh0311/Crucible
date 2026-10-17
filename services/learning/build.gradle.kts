// learning-service: M3 Lớp học, lộ trình và thi. Thiết kế tại docs/modules/m3-learning.md.
plugins {
    java
    id("org.springframework.boot") version "4.1.1"
    id("io.spring.dependency-management") version "1.1.7"
    id("org.openapi.generator") version "7.26.0"
}

group = "dev.crucible"
version = "0.1.0-SNAPSHOT"

java {
    toolchain {
        // Mặc định Java 25 (LTS); -PjavaVersion=21 cho máy chưa có JDK 25.
        languageVersion = JavaLanguageVersion.of(providers.gradleProperty("javaVersion").getOrElse("25").toInt())
    }
}

repositories {
    mavenCentral()
}

dependencies {
    implementation("dev.crucible:crucible-platform:0.1.0")
    implementation("org.springframework.boot:spring-boot-starter-actuator")
    implementation("org.springframework.boot:spring-boot-starter-amqp")
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")
    implementation("org.springframework.boot:spring-boot-starter-flyway")
    implementation("org.springframework.boot:spring-boot-starter-security")
    implementation("org.springframework.boot:spring-boot-starter-security-oauth2-resource-server")
    implementation("org.springframework.boot:spring-boot-starter-security-oauth2-client")
    implementation("org.springframework.boot:spring-boot-starter-validation")
    implementation("org.springframework.boot:spring-boot-starter-webmvc")
    implementation("org.springframework.boot:spring-boot-starter-restclient")
    implementation("org.flywaydb:flyway-database-postgresql")
    runtimeOnly("org.postgresql:postgresql")
    annotationProcessor("org.springframework.boot:spring-boot-configuration-processor")
    testImplementation("org.springframework.boot:spring-boot-starter-actuator-test")
    testImplementation("org.springframework.boot:spring-boot-starter-amqp-test")
    testImplementation("org.springframework.boot:spring-boot-starter-security-test")
    testImplementation("org.springframework.boot:spring-boot-starter-webmvc-test")
    testImplementation("org.springframework.boot:spring-boot-testcontainers")
    testImplementation("org.testcontainers:testcontainers-junit-jupiter")
    testImplementation("org.testcontainers:testcontainers-postgresql")
    testImplementation("org.testcontainers:testcontainers-rabbitmq")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher")
}

// Dạng URI để đường dẫn có dấu cách (ví dụ "Project 1") vẫn đọc được.
fun contract(name: String): String = rootDir.resolve("../../contracts/$name").canonicalFile.toURI().toString()
val generatedDir = layout.buildDirectory.get().asFile.resolve("generated/openapi")

// Client gọi content-api.yaml bằng Spring HTTP Interface, token client credentials của registration "content-service".
val generateContentClient = tasks.register<org.openapitools.generator.gradle.plugin.tasks.GenerateTask>("generateContentClient") {
    generatorName.set("spring")
    library.set("spring-http-interface")
    inputSpec.set(contract("content-api.yaml"))
    outputDir.set(generatedDir.resolve("content-client").path)
    apiPackage.set("dev.crucible.learning.infrastructure.content.generated")
    modelPackage.set("dev.crucible.learning.infrastructure.content.generated.model")
    configOptions.set(mapOf(
        "useSpringBoot4" to "true",
        "useJackson3" to "true",
        "useTags" to "true",
        "openApiNullable" to "false",
        "documentationProvider" to "none",
        "annotationLibrary" to "none",
        "configPackage" to "dev.crucible.learning.infrastructure.content.generated.config",
        "clientRegistrationId" to "content-service",
    ))
}

// Client gọi runtime-api.yaml bằng Spring HTTP Interface, token client credentials của registration "runtime-api".
val generateRuntimeClient = tasks.register<org.openapitools.generator.gradle.plugin.tasks.GenerateTask>("generateRuntimeClient") {
    generatorName.set("spring")
    library.set("spring-http-interface")
    inputSpec.set(contract("runtime-api.yaml"))
    outputDir.set(generatedDir.resolve("runtime-client").path)
    apiPackage.set("dev.crucible.learning.infrastructure.runtime.generated")
    modelPackage.set("dev.crucible.learning.infrastructure.runtime.generated.model")
    configOptions.set(mapOf(
        "useSpringBoot4" to "true",
        "useJackson3" to "true",
        "useTags" to "true",
        "openApiNullable" to "false",
        "documentationProvider" to "none",
        "annotationLibrary" to "none",
        "configPackage" to "dev.crucible.learning.infrastructure.runtime.generated.config",
        "clientRegistrationId" to "runtime-api",
    ))
}

sourceSets {
    main {
        java {
            srcDirs(generatedDir.resolve("content-client/src/main/java"), generatedDir.resolve("runtime-client/src/main/java"))
        }
    }
}

tasks.compileJava {
    dependsOn(generateContentClient, generateRuntimeClient)
}

tasks.withType<Test> {
    useJUnitPlatform()
}
