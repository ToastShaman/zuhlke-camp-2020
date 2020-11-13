package com.zuhlke.todo.client.model;

import java.util.Objects;
import java.util.StringJoiner;

public class CreateTodoResponse {
    private final Todo todo;

    public CreateTodoResponse(Todo todo) {
        this.todo = todo;
    }

    public Todo getTodo() {
        return todo;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        CreateTodoResponse that = (CreateTodoResponse) o;
        return Objects.equals(todo, that.todo);
    }

    @Override
    public int hashCode() {
        return Objects.hash(todo);
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", CreateTodoResponse.class.getSimpleName() + "[", "]")
                .add("todo=" + todo)
                .toString();
    }
}
