plugins {
    java
    id("au.com.dius.pact")
}

repositories {
    mavenCentral()
}

dependencies {
    implementation(platform("com.fasterxml.jackson:jackson-bom:_"))
    implementation("com.fasterxml.jackson.core:jackson-databind")
    implementation("com.fasterxml.jackson.module:jackson-module-parameter-names")
    implementation("com.fasterxml.jackson.datatype:jackson-datatype-jdk8")
    implementation("com.fasterxml.jackson.datatype:jackson-datatype-jsr310")

    implementation("org.slf4j:slf4j-api:_")
    implementation("ch.qos.logback:logback-classic:_")

    implementation("com.damnhandy:handy-uri-templates:_")

    testImplementation(platform(Testing.junit.bom))
    testImplementation("org.junit.jupiter:junit-jupiter")

    testImplementation("au.com.dius.pact.consumer:junit5:_")

    testImplementation("org.assertj:assertj-core:_")
    testImplementation("com.github.tomakehurst:wiremock-jre8:_")
    testImplementation("org.skyscreamer:jsonassert:_")
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