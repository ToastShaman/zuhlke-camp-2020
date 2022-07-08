package com.zuhlke.todo.client;

import com.zuhlke.todo.client.http.HttpTodoClientBuilder;
import com.zuhlke.todo.client.http.TodoClientException;
import com.zuhlke.todo.client.model.*;

public interface TodoClient {
    static HttpTodoClientBuilder builder() {
        return new HttpTodoClientBuilder();
    }

    GetTodoResponse get(GetTodoRequest request) throws TodoClientException;

    CreateTodoResponse create(CreateTodoRequest request) throws TodoClientException;

    UpdateTodoResponse update(UpdateTodoRequest request) throws TodoClientException;

    DeleteTodoResponse delete(DeleteTodoRequest request) throws TodoClientException;
}
