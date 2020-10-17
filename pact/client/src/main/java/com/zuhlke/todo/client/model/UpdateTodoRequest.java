package com.zuhlke.todo.client.model;

import com.fasterxml.jackson.annotation.JsonIgnore;

import java.util.List;
import java.util.Objects;
import java.util.StringJoiner;

public class UpdateTodoRequest {

    private final Todo todo;

    public UpdateTodoRequest(Todo todo) {
        this.todo = todo;
    }

    @JsonIgnore
    public String getId() {
        return todo.id();
    }

    public String getRev() {
        return todo.rev();
    }

    public String getText() {
        return todo.text();
    }

    public String getStatus() {
        return todo.status();
    }

    public String getCategory() {
        return todo.category();
    }

    public List<String> getTags() {
        return todo.tags();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        UpdateTodoRequest that = (UpdateTodoRequest) o;
        return Objects.equals(todo, that.todo);
    }

    @Override
    public int hashCode() {
        return Objects.hash(todo);
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", UpdateTodoRequest.class.getSimpleName() + "[", "]")
                .add("todo=" + todo)
                .toString();
    }

    public static UpdateTodoRequestBuilder builder(Todo todo) {
        return UpdateTodoRequestBuilder.builder(todo);
    }
}
