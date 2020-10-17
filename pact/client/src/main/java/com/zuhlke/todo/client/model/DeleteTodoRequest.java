package com.zuhlke.todo.client.model;

import com.fasterxml.jackson.annotation.JsonIgnore;

import java.util.Objects;
import java.util.StringJoiner;

public class DeleteTodoRequest {

    private final Todo todo;

    public DeleteTodoRequest(Todo todo) {
        this.todo = todo;
    }

    @JsonIgnore
    public String getId() {
        return todo.id();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DeleteTodoRequest that = (DeleteTodoRequest) o;
        return Objects.equals(todo, that.todo);
    }

    @Override
    public int hashCode() {
        return Objects.hash(todo);
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", DeleteTodoRequest.class.getSimpleName() + "[", "]")
                .add("todo=" + todo)
                .toString();
    }

    public static DeleteTodoRequestBuilder builder(Todo todo) {
        return DeleteTodoRequestBuilder.builder(todo);
    }
}
