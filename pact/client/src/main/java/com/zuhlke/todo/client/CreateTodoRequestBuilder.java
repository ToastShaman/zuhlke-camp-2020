package com.zuhlke.todo.client;

import java.util.List;

public final class CreateTodoRequestBuilder {
    private String text;
    private String status;
    private String category;
    private List<String> tags;

    private CreateTodoRequestBuilder() {
    }

    public static CreateTodoRequestBuilder aTodo() {
        return new CreateTodoRequestBuilder();
    }

    public CreateTodoRequestBuilder copy(CreateTodoRequest createTodoRequest) {
        this.text = createTodoRequest.getText();
        this.category = createTodoRequest.getCategory();
        this.status = createTodoRequest.getStatus();
        this.tags = List.copyOf(tags);
        return this;
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

    public CreateTodoRequestBuilder but() {
        return aTodo().withText(text).withStatus(status).withCategory(category).withTags(tags);
    }

    public CreateTodoRequest build() {
        return new CreateTodoRequest(text, status, category, tags);
    }
}
