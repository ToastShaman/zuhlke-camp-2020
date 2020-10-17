package com.zuhlke.todo.client.model;

import com.zuhlke.todo.client.TodoBuilder;

import java.util.List;
import java.util.Objects;
import java.util.StringJoiner;

public record Todo(
        String id,
        String rev,
        String text,
        String status,
        String category,
        List<String> tags
) {
    public static TodoBuilder builder() {
        return TodoBuilder.aTodo();
    }

    public static TodoBuilder newBuilder(Todo todo) {
        return TodoBuilder.aTodo(todo);
    }

    public TodoBuilder newBuilder() {
        return TodoBuilder.aTodo(this);
    }
}
