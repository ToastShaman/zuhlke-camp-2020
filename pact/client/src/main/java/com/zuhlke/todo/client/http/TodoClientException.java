package com.zuhlke.todo.client.http;

import com.zuhlke.todo.client.model.TodoApiError;

import java.util.Objects;
import java.util.Optional;
import java.util.StringJoiner;

public class TodoClientException extends RuntimeException {

    private final int status;
    private final TodoApiError error;

    public TodoClientException(int status, TodoApiError error) {
        super("Received error from Todo API");
        this.status = status;
        this.error = error;
    }

    public TodoClientException(Exception e) {
        super(e);
        this.status = -1;
        this.error = null;
    }

    public int getStatus() {
        return status;
    }

    public Optional<TodoApiError> getError() {
        return Optional.ofNullable(error);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        TodoClientException that = (TodoClientException) o;
        return status == that.status &&
                Objects.equals(error, that.error);
    }

    @Override
    public int hashCode() {
        return Objects.hash(status, error);
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", TodoClientException.class.getSimpleName() + "[", "]")
                .add("status=" + status)
                .add("error=" + error)
                .toString();
    }
}
