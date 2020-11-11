package com.zuhlke.todo.client.http;

import com.damnhandy.uri.template.UriTemplate;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.zuhlke.todo.client.*;

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.net.http.HttpResponse.BodyHandlers;
import java.time.Duration;

import static java.net.http.HttpRequest.BodyPublishers.ofByteArray;

public class HttpTodoClient implements TodoClient {

    private final String host;
    private final HttpClient client;
    private final ObjectMapper mapper;
    private final Duration timeout;

    public HttpTodoClient(String host, HttpClient client, ObjectMapper mapper, Duration timeout) {
        this.host = host;
        this.client = client;
        this.mapper = mapper;
        this.timeout = timeout;
    }

    @Override
    public CreateTodoResponse create(CreateTodoRequest createTodoRequest) {
        try {
            String uri = UriTemplate.buildFromTemplate(host)
                    .literal("/todo")
                    .build()
                    .expand();

            byte[] payload = mapper.writeValueAsBytes(createTodoRequest);

            HttpRequest request = HttpRequest.newBuilder(new URI(uri))
                    .header("Accept", "application/json")
                    .header("Content-Type", "application/json")
                    .POST(ofByteArray(payload))
                    .timeout(timeout)
                    .build();

            HttpResponse<byte[]> response = client.send(request, BodyHandlers.ofByteArray());

            if (response.statusCode() != 200 && response.statusCode() != 201) {
                TodoApiError error = mapper.readValue(response.body(), TodoApiError.class);
                throw new TodoClientException(response.statusCode(), error);
            }

            Todo todo = mapper.readValue(response.body(), Todo.class);

            return new CreateTodoResponse(todo);

        } catch (URISyntaxException | IOException | InterruptedException e) {
            throw new TodoClientException(e);
        }
    }
}
