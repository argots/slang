# Text editor interface

```slang

Screen(buffers, handle(it)).where(
  buffers: mutable(Buffers()),
  input: mutable(Buffer()),
  handle(key): case(key.string(),
    'F1': promptFileName(input, visitFile(it, buffers.get(it)))
  ),
  visitFile(name, buffer): (),
  promptFileName(buffer,





```

Ticket workflow:
1. QA creates a ticket
2. QA manager triages priority
3. Eng manager assigns to Queue
4. User picks up task
5. User completes task
5. QA verifies task
6. QA closes task or pushes it back to #4

Tables:

Untriaged Tickets
Unprioritized Tickets
User/Team Pending
User/Team InProgress
User/Team Completed
User/Team QA verified

Either done with queues or as a single tickets table + views

Declarative spec of workflow involves triggers + filters on each queue.  Hard to make sense of workflow.

Imperative flow is easier to infer intent from:

ticket = CreateTicket(Untriaged Tickets)
ticket.priority = Prioritize(ticket, priorities constraint)
ticket.queue = SelectQueue(ticket, queues constraint)
ticket.assignedTo = SelectUser(ticket, users constraint)
ticket.startTime = Start(ticket)
ticket.completedTime = Complete(ticket)
ticket.verified = Verify(ticket)

Here, each of the functions ('Prioritize', 'SelectQueue') etc can be thought of as a dialog which is "stuck" in that stage until the "workers" pick it up.  So, the task of workers is identifying which of the "ongoing workflows" they want to pick up.  The QA manager is only interested in workflows stuck at "Prioritize" for instance.  Individual users are interested in workflows stuck at "SelectUser".

Regular profiling tools and such all make sense in this context!
