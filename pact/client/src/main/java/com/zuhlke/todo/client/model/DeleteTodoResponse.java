package com.zuhlke.todo.client.model;

import java.util.Objects;
import java.util.StringJoiner;

public class DeleteTodoResponse {
    private final Todo todo;

    public DeleteTodoResponse(Todo todo) {
        this.todo = todo;
    }

    public Todo getTodo() {
        return todo;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DeleteTodoResponse that = (DeleteTodoResponse) o;
        return Objects.equals(todo, that.todo);
    }

    @Override
    public int hashCode() {
        return Objects.hash(todo);
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", DeleteTodoResponse.class.getSimpleName() + "[", "]")
                .add("todo=" + todo)
                .toString();
    }
}
