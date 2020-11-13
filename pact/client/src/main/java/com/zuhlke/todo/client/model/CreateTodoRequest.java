package com.zuhlke.todo.client.model;

import java.util.List;
import java.util.Objects;
import java.util.StringJoiner;

public class CreateTodoRequest {

    private final String text;
    private final String status;
    private final String category;
    private final List<String> tags;

    public CreateTodoRequest(String text, String status, String category, List<String> tags) {
        this.text = text;
        this.status = status;
        this.category = category;
        this.tags = List.copyOf(tags);
    }

    public String getText() {
        return text;
    }

    public String getStatus() {
        return status;
    }

    public String getCategory() {
        return category;
    }

    public List<String> getTags() {
        return tags;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        CreateTodoRequest createTodoRequest = (CreateTodoRequest) o;
        return Objects.equals(text, createTodoRequest.text) &&
                Objects.equals(status, createTodoRequest.status) &&
                Objects.equals(category, createTodoRequest.category) &&
                Objects.equals(tags, createTodoRequest.tags);
    }

    @Override
    public int hashCode() {
        return Objects.hash(text, status, category, tags);
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", CreateTodoRequest.class.getSimpleName() + "[", "]")
                .add("text='" + text + "'")
                .add("status='" + status + "'")
                .add("category='" + category + "'")
                .add("tags=" + tags)
                .toString();
    }

    public static CreateTodoRequestBuilder builder() {
        return CreateTodoRequestBuilder.builder();
    }
}
