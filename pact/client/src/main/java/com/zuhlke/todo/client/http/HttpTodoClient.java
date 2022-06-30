package com.zuhlke.todo.client.http;

import com.damnhandy.uri.template.UriTemplate;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.zuhlke.todo.client.*;
import com.zuhlke.todo.client.model.*;

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.net.http.HttpResponse.BodyHandlers;
import java.time.Duration;
import java.util.Objects;
import java.util.Set;

import static com.damnhandy.uri.template.UriTemplateBuilder.var;
import static java.net.http.HttpRequest.BodyPublishers.ofByteArray;

@SuppressWarnings("DuplicatedCode")
public class HttpTodoClient implements TodoClient {

    private final String host;
    private final HttpClient client;
    private final ObjectMapper mapper;
    private final Duration timeout;

    public HttpTodoClient(String host,
                          HttpClient client,
                          ObjectMapper mapper,
                          Duration timeout) {
        this.host = host;
        this.client = client;
        this.mapper = mapper;
        this.timeout = timeout;
    }

    @Override
    public GetTodoResponse get(GetTodoRequest getTodoRequest) {
        try {
            String uri = UriTemplate.buildFromTemplate(host)
                    .literal("/todo")
                    .path(var("id"))
                    .build()
                    .set("id", getTodoRequest.id())
                    .expand();

            HttpRequest request = HttpRequest.newBuilder(new URI(uri))
                    .header("Accept", "application/json; charset=UTF-8")
                    .GET()
                    .timeout(timeout)
                    .build();

            HttpResponse<byte[]> response = client.send(request, BodyHandlers.ofByteArray());

            if (!Objects.equals(200, response.statusCode())) {
                TodoApiError error = mapper.readValue(response.body(), TodoApiError.class);
                throw new TodoClientException(response.statusCode(), error);
            }

            Todo todo = mapper.readValue(response.body(), Todo.class);

            return new GetTodoResponse(todo);

        } catch (URISyntaxException | IOException | InterruptedException e) {
            throw new TodoClientException(e);
        }
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
                    .header("Accept", "application/json; charset=UTF-8")
                    .header("Content-Type", "application/json; charset=UTF-8")
                    .POST(ofByteArray(payload))
                    .timeout(timeout)
                    .build();

            HttpResponse<byte[]> response = client.send(request, BodyHandlers.ofByteArray());

            if (!Set.of(200, 201).contains(response.statusCode())) {
                TodoApiError error = mapper.readValue(response.body(), TodoApiError.class);
                throw new TodoClientException(response.statusCode(), error);
            }

            Todo todo = mapper.readValue(response.body(), Todo.class);

            return new CreateTodoResponse(todo);

        } catch (URISyntaxException | IOException | InterruptedException e) {
            throw new TodoClientException(e);
        }
    }

    @Override
    public UpdateTodoResponse update(UpdateTodoRequest updateTodoRequest) {
        try {
            String uri = UriTemplate.buildFromTemplate(host)
                    .literal("/todo")
                    .path(var("id"))
                    .build()
                    .set("id", updateTodoRequest.getId())
                    .expand();

            byte[] payload = mapper.writeValueAsBytes(updateTodoRequest);

            HttpRequest request = HttpRequest.newBuilder(new URI(uri))
                    .header("Accept", "application/json; charset=UTF-8")
                    .header("Content-Type", "application/json; charset=UTF-8")
                    .PUT(ofByteArray(payload))
                    .timeout(timeout)
                    .build();

            HttpResponse<byte[]> response = client.send(request, BodyHandlers.ofByteArray());

            if (!Set.of(200, 201).contains(response.statusCode())) {
                TodoApiError error = mapper.readValue(response.body(), TodoApiError.class);
                throw new TodoClientException(response.statusCode(), error);
            }

            Todo todo = mapper.readValue(response.body(), Todo.class);

            return new UpdateTodoResponse(todo);

        } catch (URISyntaxException | IOException | InterruptedException e) {
            throw new TodoClientException(e);
        }
    }

    @Override
    public DeleteTodoResponse delete(DeleteTodoRequest deleteTodoRequest) {
        try {
            String uri = UriTemplate.buildFromTemplate(host)
                    .literal("/todo")
                    .path(var("id"))
                    .build()
                    .set("id", deleteTodoRequest.getId())
                    .expand();

            HttpRequest request = HttpRequest.newBuilder(new URI(uri))
                    .header("Accept", "application/json; charset=UTF-8")
                    .DELETE()
                    .timeout(timeout)
                    .build();

            HttpResponse<byte[]> response = client.send(request, BodyHandlers.ofByteArray());

            if (response.statusCode() != 200) {
                TodoApiError error = mapper.readValue(response.body(), TodoApiError.class);
                throw new TodoClientException(response.statusCode(), error);
            }

            Todo todo = mapper.readValue(response.body(), Todo.class);

            return new DeleteTodoResponse(todo);

        } catch (URISyntaxException | IOException | InterruptedException e) {
            throw new TodoClientException(e);
        }
    }
}
