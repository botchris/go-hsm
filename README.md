# UML HSM

HSM package provides a simple state chart library written in Go.

## Supported UML State Chart Features

| Feature              | Implemented | Test case            |
|----------------------|:-----------:|----------------------|
| Simple state         |     Yes     | door_test            |
| Composite states     |     Yes     | nesting_test         |
| Sub machines         |     No      |                      |
| Compound transition  |     No      |                      |
| Fork                 |     No      |                      |
| Join                 |     No      |                      |
| Guards/Actions       |     Yes     | lobby_test + various |
| Shallow/Deep history |     No      |                      | 
| Exit/Entry points    |     Yes     | error_test + various |
| Init/Final           |     Yes     | various              |
| Event deferral       |     No      |                      |
| Terminate            |     No      |                      |
| Choice               |     Yes     | choice_test          |
| Junction             |     No      |                      |
| Do activity          |     No      |                      |
| Connection point ref |     No      |                      |
| Protocol Machines    |     No      |                      |

## Introduction

State machines are used to model the dynamic behavior of a model element, and more specifically, the event-driven
aspects of the system's behavior (see Events and Signals). State machines are specifically used to define
state-dependent behavior, or behavior that varies depending on the state in which the model element is in. Model
elements whose behavior does not vary with its state of the element do not require state machines to describe their
behavior (these elements are typically passive classes whose primary responsible is to manage data). In particular,
state machines must be used to model the behavior of active classes that use call events and signal events to implement
their operations (as transitions in the class's state machine).

Source: [Guideline: Statechart Diagram](http://www.michael-richardson.com/processes/rup_for_sqa)

## States

A state is a condition of an object in which it performs some task or waits for an event. An object may remain in a
state for a finite amount of time. A state has several properties:

- **Name**: A textual string which distinguishes the state from other states.
- **Entry/exit actions**: Actions executed on entering and exiting the state.
- **Internal transitions**: Transitions that are handled without causing a change in state.
- **Sub-states**: The nested structure of a state.
- **Deferred events**: TODO

### Choice Pseudo-States

Realizes a dynamic conditional branch. It evaluates the guards of the triggers of its outgoing transitions to select
only one outgoing transition.

#### Entry and Exit Actions

Entry and exit actions allow the same action to be dispatched every time the state is entered or left, respectively.
Entry and exit actions enable this to be done cleanly, without having to explicitly put the actions on every incoming or
outgoing transition explicitly.

## Transitions

A transition is a relationship between two states indicating that an object in the first state will perform certain
actions and enter a second state when a specified event occurs and specified conditions are satisfied. On such a change
of state, the transition is said to 'fire'. Until the transition fires, the object is said to be in the 'source' state;
after it fires, it is said to be in the 'target' state. A transition has several properties:

- **Source state**: The state affected by the transition; if an object is in the source state, an outgoing transition
  may fire when the object receives the trigger event of the transition and if the guard condition, if any, is satisfied.
- **Event trigger**: The event that makes the transition eligible to fire (providing its guard condition is satisfied)
  when received by the object in the source state.
- **Guard condition**: A boolean expression that is evaluated when the transition is triggered by the reception of the
  event trigger; if the expression evaluates True, the transition is eligible to fire; if the expression evaluates to
  False, the transition does not fire. If there is no other transition that could be triggered by the same event, the
  event is lost.
- **Effect**: An executable atomic computation that may directly act upon the object that owns the state machine, and
  indirectly on other objects that are visible to the object.
- **Target state**: The state that is active after the completion of the transition.

### Signals

In the context of the state machine, a Signal is an occurrence of a stimulus that can trigger a state transition.
Signals may include the passing of time, or a change in state. A signal or call may have parameters whose values are
available to the transition, including expressions for the guard conditions and action. It is also possible to have a
signal-less transition, represented by a transition with no signal trigger. These transitions, also called completion
transitions, are triggered implicitly when its source state has completed its actions.

### Guards

Transition guard conditions are evaluated after the signal for the transition occurs. It is possible to have multiple
transitions from the same source state and with the same signal trigger, as long as the guard conditions don't overlap.
A guard condition is evaluated just once for the transition at the time the signal occurs. The boolean MUST be
side effect free, at least none that would alter evaluation of other guards having the same trigger.

### Effects

A transition effect is an executable atomic computation, meaning that it cannot be interrupted by an event and therefore
runs to completion. Effects may include operation calls (to the owner of the state machine as well as other visible
objects), the creation or destruction of another object, or the sending of a signal to another object.

### Internal Transitions

Are those that may have an effect but not a change of state. Internal transitions allow signals to be handled within the
state without leaving the state, thereby avoiding triggering entry or exit actions. Internal transitions may have guard
conditions, and essentially represent interrupt-handlers.

# Concepts

**Events and Signals**

An event is the specification of a significant occurrence that has a location in time and space. A 'signal' is a kind of
event that represents the specification of an asynchronous stimulus between two instances.