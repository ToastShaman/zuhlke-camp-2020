package com.zuhlke.todo.client.model;

import java.util.List;

public final class CreateTodoRequestBuilder {
    private String text;
    private String status;
    private String category;
    private List<String> tags;

    private CreateTodoRequestBuilder() {
    }

    public static CreateTodoRequestBuilder builder() {
        return new CreateTodoRequestBuilder();
    }

    public CreateTodoRequestBuilder withText(String text) {
        this.text = text;
        return this;
    }

    public CreateTodoRequestBuilder withStatus(String status) {
        this.status = status;
        return this;
    }

    public CreateTodoRequestBuilder withCategory(String category) {
        this.category = category;
        return this;
    }

    public CreateTodoRequestBuilder withTags(List<String> tags) {
        this.tags = tags;
        return this;
    }

    public CreateTodoRequest build() {
        return new CreateTodoRequest(text, status, category, tags);
    }
}
