package com.zuhlke.todo.client.model;

import java.util.List;
import java.util.Objects;
import java.util.StringJoiner;

public record TodoApiError(String message, List<String> details) {

    public TodoApiError(String message, List<String> details) {
        this.message = message;
        this.details = details == null ? List.of() : List.copyOf(details);
    }
}
