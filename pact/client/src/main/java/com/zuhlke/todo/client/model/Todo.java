package com.zuhlke.todo.client.model;

import com.zuhlke.todo.client.TodoBuilder;

import java.util.List;
import java.util.Objects;
import java.util.StringJoiner;

public class Todo {

    private final String id;
    private final String rev;
    private final String text;
    private final String status;
    private final String category;
    private final List<String> tags;

    public Todo(String id,
                String rev,
                String text,
                String status,
                String category,
                List<String> tags) {
        this.id = id;
        this.rev = rev;
        this.text = text;
        this.status = status;
        this.category = category;
        this.tags = List.copyOf(tags);
    }

    public String getId() {
        return id;
    }

    public String getRev() {
        return rev;
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
        Todo todo = (Todo) o;
        return Objects.equals(id, todo.id) &&
                Objects.equals(rev, todo.rev) &&
                Objects.equals(text, todo.text) &&
                Objects.equals(status, todo.status) &&
                Objects.equals(category, todo.category) &&
                Objects.equals(tags, todo.tags);
    }

    @Override
    public int hashCode() {
        return Objects.hash(id, rev, text, status, category, tags);
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", Todo.class.getSimpleName() + "[", "]")
                .add("id='" + id + "'")
                .add("rev='" + rev + "'")
                .add("text='" + text + "'")
                .add("status='" + status + "'")
                .add("category='" + category + "'")
                .add("tags=" + tags)
                .toString();
    }

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
