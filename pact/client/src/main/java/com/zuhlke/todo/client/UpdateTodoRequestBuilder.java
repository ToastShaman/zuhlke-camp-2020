package com.zuhlke.todo.client;

import java.util.List;

public class UpdateTodoRequestBuilder {

    private final TodoBuilder todoBuilder;

    private UpdateTodoRequestBuilder(TodoBuilder todoBuilder) {
        this.todoBuilder = todoBuilder;
    }

    public static UpdateTodoRequestBuilder aUpdateTodoRequestBuilder(Todo todo) {
        return new UpdateTodoRequestBuilder(todo.newBuilder());
    }

    public UpdateTodoRequestBuilder withId(String id) {
        todoBuilder.withId(id);
        return this;
    }

    public UpdateTodoRequestBuilder withRev(String rev) {
        todoBuilder.withRev(rev);
        return this;
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
