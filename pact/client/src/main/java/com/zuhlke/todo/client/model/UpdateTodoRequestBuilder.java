package com.zuhlke.todo.client.model;

import com.zuhlke.todo.client.TodoBuilder;

import java.util.List;

public class UpdateTodoRequestBuilder {

    private final TodoBuilder todoBuilder;

    private UpdateTodoRequestBuilder(TodoBuilder todoBuilder) {
        this.todoBuilder = todoBuilder;
    }

    public static UpdateTodoRequestBuilder builder(Todo todo) {
        return new UpdateTodoRequestBuilder(todo.newBuilder());
    }

    public UpdateTodoRequestBuilder withText(String text) {
        todoBuilder.withText(text);
        return this;
    }

    public UpdateTodoRequestBuilder withStatus(String status) {
        todoBuilder.withStatus(status);
        return this;
    }

    public UpdateTodoRequestBuilder withCategory(String category) {
        todoBuilder.withCategory(category);
        return this;
    }

    public UpdateTodoRequestBuilder withTags(List<String> tags) {
        todoBuilder.withTags(tags);
        return this;
    }

    public UpdateTodoRequest build() {
        return new UpdateTodoRequest(todoBuilder.build());
    }
}
