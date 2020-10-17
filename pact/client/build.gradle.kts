plugins {
    java
    id("au.com.dius.pact") version "4.1.0"
}

repositories {
    mavenCentral()
}

dependencies {
    implementation(platform("com.fasterxml.jackson:jackson-bom:2.13.3"))
    implementation("com.fasterxml.jackson.core:jackson-databind")
    implementation("com.fasterxml.jackson.module:jackson-module-parameter-names")
    implementation("com.fasterxml.jackson.datatype:jackson-datatype-jdk8")
    implementation("com.fasterxml.jackson.datatype:jackson-datatype-jsr310")

    implementation("org.slf4j:slf4j-api:1.7.36")
    implementation("ch.qos.logback:logback-classic:1.2.11")

    implementation("com.damnhandy:handy-uri-templates:2.1.8")

    testImplementation(platform("org.junit:junit-bom:5.8.2"))
    testImplementation("org.junit.jupiter:junit-jupiter")

    testImplementation("au.com.dius.pact.consumer:junit5:4.3.9")

    testImplementation("org.assertj:assertj-core:3.23.1")
    testImplementation("com.github.tomakehurst:wiremock-jre8:2.33.2")
    testImplementation("org.skyscreamer:jsonassert:1.5.0")
}

group = "com.zuhlke"
version = "0.0.1-SNAPSHOT"
java.sourceCompatibility = JavaVersion.VERSION_18
java.targetCompatibility = JavaVersion.VERSION_18

tasks.test {
    useJUnitPlatform()
    testLogging {
        events("passed", "skipped", "failed")
    }
}

pact {
    publish {
        pactBrokerUrl = "http://localhost:9292"
    }
    reports {
        defaultReports()
    }
    serviceProviders {
        create("provider1") {
            protocol = "http"
        }
    }
}