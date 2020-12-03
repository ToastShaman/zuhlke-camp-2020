package com.zuhlke.todo.client;

import com.zuhlke.todo.client.http.HttpTodoClientBuilder;
import com.zuhlke.todo.client.model.*;

public interface TodoClient {
    static HttpTodoClientBuilder builder() {
        return new HttpTodoClientBuilder();
    }

    GetTodoResponse get(GetTodoRequest request);

    CreateTodoResponse create(CreateTodoRequest request);

    UpdateTodoResponse update(UpdateTodoRequest request);

    DeleteTodoResponse delete(DeleteTodoRequest request);
}
