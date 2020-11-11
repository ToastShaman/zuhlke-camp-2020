package com.zuhlke.todo.client;

import com.zuhlke.todo.client.http.HttpTodoClientBuilder;

public interface TodoClient {
    static HttpTodoClientBuilder builder() {
        return new HttpTodoClientBuilder();
    }

    CreateTodoResponse create(CreateTodoRequest createTodoRequest);

    UpdateTodoResponse update(UpdateTodoRequest updateTodoRequest);
}
