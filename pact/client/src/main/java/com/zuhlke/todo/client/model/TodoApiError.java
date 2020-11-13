package com.zuhlke.todo.client.model;

import java.util.List;
import java.util.Objects;
import java.util.StringJoiner;

public class TodoApiError {

    private final String message;
    private final List<String> details;

    public TodoApiError(String message, List<String> details) {
        this.message = message;
        this.details = details == null ? List.of() : List.copyOf(details);
    }

    public List<String> getDetails() {
        return details;
    }

    public String getMessage() {
        return message;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        TodoApiError that = (TodoApiError) o;
        return Objects.equals(message, that.message) &&
                Objects.equals(details, that.details);
    }

    @Override
    public int hashCode() {
        return Objects.hash(message, details);
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", TodoApiError.class.getSimpleName() + "[", "]")
                .add("message='" + message + "'")
                .add("details=" + details)
                .toString();
    }
}
