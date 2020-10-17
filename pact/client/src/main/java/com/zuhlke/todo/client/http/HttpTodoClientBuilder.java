package com.zuhlke.todo.client.http;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.zuhlke.todo.client.SystemObjectMapper;

import java.net.http.HttpClient;
import java.time.Duration;
import java.util.Objects;

public class HttpTodoClientBuilder {

    private ObjectMapper objectMapper = SystemObjectMapper.createDefaultMapper();
    private HttpClient httpClient = HttpClient.newHttpClient();
    private Duration timeout = Duration.ofSeconds(3);
    private String host;

    public HttpTodoClientBuilder setObjectMapper(ObjectMapper objectMapper) {
        this.objectMapper = objectMapper;
        return this;
    }

    public HttpTodoClientBuilder setHttpClient(HttpClient httpClient) {
        this.httpClient = httpClient;
        return this;
    }

    public HttpTodoClientBuilder setTimeout(Duration timeout) {
        this.timeout = timeout;
        return this;
    }

    public HttpTodoClientBuilder setHost(String host) {
        this.host = host;
        return this;
    }

    public HttpTodoClient build() {
        return new HttpTodoClient(
                Objects.requireNonNull(host),
                Objects.requireNonNull(httpClient),
                Objects.requireNonNull(objectMapper),
                Objects.requireNonNull(timeout)
        );
    }
}
