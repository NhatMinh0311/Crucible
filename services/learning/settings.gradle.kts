plugins {
    id("org.gradle.toolchains.foojay-resolver-convention") version "1.0.0"
}

rootProject.name = "learning-service"

// Thư viện dùng chung của nền tảng (outbox, inbox, xác thực).
includeBuild("../../libs/java/crucible-platform")
