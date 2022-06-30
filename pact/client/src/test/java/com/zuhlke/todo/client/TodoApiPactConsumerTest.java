package com.zuhlke.todo.client;

import au.com.dius.pact.consumer.MockServer;
import au.com.dius.pact.consumer.dsl.PactDslJsonBody;
import au.com.dius.pact.consumer.dsl.PactDslWithProvider;
import au.com.dius.pact.consumer.junit5.PactConsumerTestExt;
import au.com.dius.pact.consumer.junit5.PactTestFor;
import au.com.dius.pact.core.model.PactSpecVersion;
import au.com.dius.pact.core.model.RequestResponsePact;
import au.com.dius.pact.core.model.annotations.Pact;
import com.zuhlke.todo.client.http.HttpTodoClient;
import com.zuhlke.todo.client.model.*;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;

import java.util.List;

import static org.assertj.core.api.AssertionsForClassTypes.assertThat;

//@formatter:off
@SuppressWarnings({"unused", "ConstantConditions"})
@ExtendWith(PactConsumerTestExt.class)
@PactTestFor(providerName = "todo_api", port = "1234", pactVersion = PactSpecVersion.V3)
public class TodoApiPactConsumerTest {

    @Pact(consumer = "jvm_todo_client")
    public RequestResponsePact createTodo(PactDslWithProvider builder) {
        return builder
                .given("an empty repository")
                .uponReceiving("a new todo")
                    .method("POST")
                    .path("/todo")
                    .body("""
                        {
                            "status": "TODO",
                            "text": "Don't forget the milk",
                            "category": "shopping",
                            "tags": [
                                "groceries",
                                "food"
                            ]
                        }
                        """, "application/json")
                    .matchHeader("Content-Type", "application/json", "application/json")
                .willRespondWith()
                    .status(201)
                    .body(new PactDslJsonBody()
                            .stringType("id")
                            .stringType("rev")
                            .stringType("text")
                            .stringType("status")
                            .stringType("category")
                            .array("tags")
                                .stringType()
                                .stringType()
                            .closeArray())
                    .matchHeader("Content-Type", "application/json", "application/json")
                .toPact();
    }

    @Test
    @PactTestFor(pactMethod = "createTodo")
    void testCreateTodo(MockServer mockServer)  {
        HttpTodoClient todoClient = TodoClient.builder()
                .setHost(mockServer.getUrl())
                .build();

        CreateTodoRequest request = CreateTodoRequest.builder()
                .withText("Don't forget the milk")
                .withStatus("TODO")
                .withCategory("shopping")
                .withTags(List.of("groceries", "food"))
                .build();

        CreateTodoResponse response = todoClient.create(request);

        Todo actual = response.todo();

        assertThat(actual).hasNoNullFieldsOrProperties();
    }

    @Pact(consumer = "jvm_todo_client")
    public RequestResponsePact updateTodo(PactDslWithProvider builder) {
        return builder
                .given("an existing todo with id=--MLqrG6LkLkkKc1iMLBt")
                .uponReceiving("an update")
                .method("PUT")
                .path("/todo/-MLqrG6LkLkkKc1iMLBt")
                .body("""
                        {
                            "rev": "-MLivp1BrS59mMbSN7Jr",
                            "status": "DONE",
                            "text": "Don't forget the milk",
                            "category": "shopping",
                            "tags": [
                                "groceries",
                                "food"
                            ]
                        }
                        """, "application/json")
                .matchHeader("Content-Type", "application/json", "application/json")
                .willRespondWith()
                .status(200)
                .body(new PactDslJsonBody()
                        .stringType("id")
                        .stringType("rev")
                        .stringType("text")
                        .stringType("status")
                        .stringType("category")
                        .array("tags")
                            .stringType()
                            .stringType()
                        .closeArray())
                .matchHeader("Content-Type", "application/json", "application/json")
                .toPact();
    }

    @Test
    @PactTestFor(pactMethod = "updateTodo")
    void testUpdateTodo(MockServer mockServer)  {
        HttpTodoClient todoClient = TodoClient.builder()
                .setHost(mockServer.getUrl())
                .build();

        Todo existing = Todo.builder()
                .withId("-MLqrG6LkLkkKc1iMLBt")
                .withRev("-MLivp1BrS59mMbSN7Jr")
                .withText("Don't forget the milk")
                .withStatus("TODO")
                .withCategory("shopping")
                .withTags(List.of("groceries", "food"))
                .build();

        UpdateTodoRequest request = UpdateTodoRequest.builder(existing)
                .withText("Don't forget the milk")
                .withStatus("DONE")
                .withCategory("shopping")
                .withTags(List.of("groceries", "food"))
                .build();

        Todo actual = todoClient.update(request).todo();

        assertThat(actual).hasNoNullFieldsOrProperties();
    }

    @Pact(consumer = "jvm_todo_client")
    public RequestResponsePact deleteTodo(PactDslWithProvider builder) {
        return builder
                .given("an existing todo with id=--MLqrG6LkLkkKc1iMLBt")
                .uponReceiving("a deletion")
                .method("DELETE")
                .path("/todo/-MLqrG6LkLkkKc1iMLBt")
                .willRespondWith()
                .status(200)
                .body(new PactDslJsonBody()
                        .stringType("id")
                        .stringType("rev")
                        .stringType("text")
                        .stringType("status")
                        .stringType("category")
                        .array("tags")
                            .stringType()
                            .stringType()
                        .closeArray())
                .matchHeader("Content-Type", "application/json", "application/json")
                .toPact();
    }

    @Test
    @PactTestFor(pactMethod = "deleteTodo")
    void testDeleteTodo(MockServer mockServer)  {
        HttpTodoClient todoClient = TodoClient.builder()
                .setHost(mockServer.getUrl())
                .build();

        Todo existing = Todo.builder()
                .withId("-MLqrG6LkLkkKc1iMLBt")
                .withRev("-MLivp1BrS59mMbSN7Jr")
                .withText("Don't forget the milk")
                .withStatus("TODO")
                .withCategory("shopping")
                .withTags(List.of("groceries", "food"))
                .build();

        DeleteTodoRequest request = DeleteTodoRequest.builder(existing).build();

        Todo actual = todoClient.delete(request).todo();

        assertThat(actual).hasNoNullFieldsOrProperties();
    }

    @Pact(consumer = "jvm_todo_client")
    public RequestResponsePact getTodo(PactDslWithProvider builder) {
        return builder
                .given("an existing todo with id=--MLqrG6LkLkkKc1iMLBt")
                .uponReceiving("a retrieval")
                .method("GET")
                .path("/todo/-MLqrG6LkLkkKc1iMLBt")
                .willRespondWith()
                .status(200)
                .body(new PactDslJsonBody()
                        .stringType("id")
                        .stringType("rev")
                        .stringType("text")
                        .stringType("status")
                        .stringType("category")
                        .array("tags")
                            .stringType()
                            .stringType()
                        .closeArray())
                .matchHeader("Content-Type", "application/json", "application/json")
                .toPact();
    }

    @Test
    @PactTestFor(pactMethod = "getTodo")
    void testGetTodo(MockServer mockServer)  {
        HttpTodoClient todoClient = TodoClient.builder()
                .setHost(mockServer.getUrl())
                .build();

        Todo existing = Todo.builder()
                .withId("-MLqrG6LkLkkKc1iMLBt")
                .withRev("-MLivp1BrS59mMbSN7Jr")
                .withText("Don't forget the milk")
                .withStatus("TODO")
                .withCategory("shopping")
                .withTags(List.of("groceries", "food"))
                .build();

        GetTodoRequest request = new GetTodoRequest("-MLqrG6LkLkkKc1iMLBt");

        Todo todo = todoClient.get(request).todo();

        assertThat(todo).hasNoNullFieldsOrProperties();
    }
}
