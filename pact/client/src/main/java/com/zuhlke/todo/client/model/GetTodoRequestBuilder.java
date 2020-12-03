package com.zuhlke.todo.client.model;

public final class GetTodoRequestBuilder {
    private String id;

    private GetTodoRequestBuilder() {
    }

    public static GetTodoRequestBuilder builder() {
        return new GetTodoRequestBuilder();
    }

    public GetTodoRequestBuilder withId(String id) {
        this.id = id;
        return this;
    }

    public GetTodoRequest build() {
        return new GetTodoRequest(id);
    }
}
