package com.zuhlke.todo.client;

import au.com.dius.pact.consumer.MockServer;
import au.com.dius.pact.consumer.dsl.PactDslJsonBody;
import au.com.dius.pact.consumer.dsl.PactDslWithProvider;
import au.com.dius.pact.consumer.junit5.PactConsumerTestExt;
import au.com.dius.pact.consumer.junit5.PactTestFor;
import au.com.dius.pact.core.model.RequestResponsePact;
import au.com.dius.pact.core.model.annotations.Pact;
import com.zuhlke.todo.client.http.HttpTodoClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;

import java.util.List;

import static org.assertj.core.api.AssertionsForClassTypes.assertThat;

//@formatter:off
@ExtendWith(PactConsumerTestExt.class)
@PactTestFor(providerName = "todo_api", port = "1234")
public class TodoApiPactConsumerTest {

    @BeforeEach
    public void setUp(MockServer mockServer) {
        assertThat(mockServer).isNotNull();
    }

    @Pact(consumer = "jvm_todo_client")
    public RequestResponsePact createTodo(PactDslWithProvider builder) {
        return builder
                .given("An empty repository")
                .uponReceiving("a new todo")
                    .method("POST")
                    .path("/todo")
                    .body("""
                        {
                            "text": "Don't forget the milk",
                            "status": "TODO",
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
                            .stringMatcher("id", "[-a-zA-Z0-9_]{20}", "-MLqrG6LkLkkKc1iMLBt")
                            .stringMatcher("rev", "[-a-zA-Z0-9_]{20}", "-MLivp1BrS59mMbSN7Jr")
                            .stringValue("text", "Don't forget the milk")
                            .stringValue("status", "TODO")
                            .stringValue("category", "shopping")
                            .array("tags").includesStr("groceries").includesStr("food").closeArray())
                    .matchHeader("Content-Type", "application/json", "application/json")
                .toPact();
    }

    @Test
    @PactTestFor(pactMethod = "createTodo")
    void testCreateTodo(MockServer mockServer)  {
        HttpTodoClient todoClient = TodoClient.builder()
                .setHost(mockServer.getUrl())
                .build();

        CreateTodoRequest createTodoRequest = CreateTodoRequest.builder()
                .withText("Don't forget the milk")
                .withStatus("TODO")
                .withCategory("shopping")
                .withTags(List.of("groceries", "food"))
                .build();

        CreateTodoResponse response = todoClient.create(createTodoRequest);

        Todo actual = response.getTodo();

        Todo expected = Todo.builder()
                .withId("-MLqrG6LkLkkKc1iMLBt")
                .withRev("-MLivp1BrS59mMbSN7Jr")
                .withText("Don't forget the milk")
                .withStatus("TODO")
                .withCategory("shopping")
                .withTags(List.of("groceries", "food"))
                .build();

        assertThat(actual).isEqualTo(expected);
    }
}
