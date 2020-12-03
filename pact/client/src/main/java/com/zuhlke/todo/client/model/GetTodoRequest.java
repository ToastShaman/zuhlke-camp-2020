package com.zuhlke.todo.client.model;

import java.util.Objects;
import java.util.StringJoiner;

public class GetTodoRequest {

    private final String id;

    public GetTodoRequest(String id) {
        this.id = id;
    }

    public String getId() {
        return id;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        GetTodoRequest that = (GetTodoRequest) o;
        return Objects.equals(id, that.id);
    }

    @Override
    public int hashCode() {
        return Objects.hash(id);
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", GetTodoRequest.class.getSimpleName() + "[", "]")
                .add("id='" + id + "'")
                .toString();
    }
}
