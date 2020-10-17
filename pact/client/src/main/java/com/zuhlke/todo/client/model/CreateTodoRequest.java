package com.zuhlke.todo.client.model;

import java.util.List;

public record CreateTodoRequest(
        String text,
        String status,
        String category,
        List<String> tags
) {
    public static CreateTodoRequestBuilder builder() {
        return CreateTodoRequestBuilder.builder();
    }
}
