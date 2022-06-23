# zuhlke-camp-2020

Consumer Driven Contract Testing demo using [Pact](https://docs.pact.io/)

## Prerequisites

1. Find a pair to do the exercise with

1. If you are doing this exercise remotely:

    1. Each of you needs to install [Tuple][4] and register with your **Zuhlke email address**!

    1. Ask Kevin for the magic team link

    1. Start a call

1. Read the [Etiquette for Pair Programming][1]

1. Choose a technology for the *Consumer* you are going to write

1. Choose a technology for the *Provider* you are going to write

## Instructions

You'll be practising writing API contracts by using [PACT][2] in the language/tech-stack of your choice.

The *Provider* adheres to the following [API Specification](#api-specifications).

1. Start implementing the *Consumer* **incrementally** by using [PACT][2] and once you have a working test/contract for one of the specifications share the contract with the *Provider*.

1. Run the contract you just been given from the *Consumer* against the *Provider* and try to make it pass.

1. Repeat the process as you implement more and more of the specification.

## API Specification

Copied from [TODO-Backend API][3]

### the pre-requisites

* the api root responds to a GET (i.e. the server is up and accessible, CORS headers are set up)
* the api root responds to a POST with the todo which was posted to it
* the api root responds successfully to a DELETE
* after a DELETE the api root responds to a GET with a JSON representation of an empty array

### storing new todos by posting to the root url

* adds a new todo to the list of todos at the root url
* sets up a new todo as initially not completed
* each new todo has a url
* each new todo has a url, which returns a todo

### working with an existing todo

* can navigate from a list of todos to an individual todo via urls
* can change the todo's title by PATCHing to the todo's url
* can change the todo's completedness by PATCHing to the todo's url
* changes to a todo are persisted and show up when re-fetching the todo
* can delete a todo making a DELETE request to the todo's url

### tracking todo order

* can create a todo with an order field
* can PATCH a todo to change its order
* remembers changes to a todo's order

[1]: https://github.com/ToastShaman/zuhlke-camp-2022/blob/open-space/pair-programming/build/pair-programming-guide.pdf
[2]: https://pact.io/
[3]: https://todobackend.com/
[4]: https://tuple.app/