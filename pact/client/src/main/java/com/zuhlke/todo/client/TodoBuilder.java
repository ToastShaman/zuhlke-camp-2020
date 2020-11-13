package com.zuhlke.todo.client;

import com.zuhlke.todo.client.model.Todo;

import java.util.List;

public final class TodoBuilder {
    private String id;
    private String rev;
    private String text;
    private String status;
    private String category;
    private List<String> tags;

    private TodoBuilder() {
    }

    public static TodoBuilder aTodo() {
        return new TodoBuilder();
    }

    public static TodoBuilder aTodo(Todo todo) {
        return new TodoBuilder()
                .withId(todo.getId())
                .withRev(todo.getRev())
                .withStatus(todo.getStatus())
                .withCategory(todo.getCategory())
                .withTags(todo.getTags())
                .withText(todo.getText());
    }

    public TodoBuilder withId(String id) {
        this.id = id;
        return this;
    }

    public TodoBuilder withRev(String rev) {
        this.rev = rev;
        return this;
    }

    public TodoBuilder withText(String text) {
        this.text = text;
        return this;
    }

    public TodoBuilder withStatus(String status) {
        this.status = status;
        return this;
    }

    public TodoBuilder withCategory(String category) {
        this.category = category;
        return this;
    }

    public TodoBuilder withTags(List<String> tags) {
        this.tags = tags;
        return this;
    }

    public Todo build() {
        return new Todo(id, rev, text, status, category, tags);
    }
}
