package com.zuhlke.todo.client.model;

import com.zuhlke.todo.client.TodoBuilder;

import java.util.List;

public class DeleteTodoRequestBuilder {

    private final TodoBuilder todoBuilder;

    private DeleteTodoRequestBuilder(TodoBuilder todoBuilder) {
        this.todoBuilder = todoBuilder;
    }

    public static DeleteTodoRequestBuilder builder(Todo todo) {
        return new DeleteTodoRequestBuilder(todo.newBuilder());
    }

    public DeleteTodoRequest build() {
        return new DeleteTodoRequest(todoBuilder.build());
    }
}
