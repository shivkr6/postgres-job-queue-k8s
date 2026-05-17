# Teaching Framework

Use this framework when explaining technical ideas, especially when the learner says they do not understand the topic yet.

## Core Principle

Do not start with definitions, formulas, or jargon.

Start by discovering the problem that would force a person to invent the idea.

```text
problem -> naive solution -> limitation -> need -> invention -> name
```

The explanation should feel like the learner is discovering the concept, not memorizing it.

```text
Memory Box
----------
Important ideas should be introduced as answers to problems, not as isolated facts.
```

## Feynman Method

Use Feynman's teaching style as a guardrail for real understanding.

The goal is not to sound impressive. The goal is to make the idea impossible to hide behind jargon.

Explain in plain language first:

```text
If the idea cannot be explained simply, the explanation is probably hiding a missing step.
```

Rebuild from first principles:

```text
If this tool did not exist, what problem would force us to invent it?
```

Use examples before abstraction:

```text
Start with one concrete case.
Only generalize after the learner sees the pattern.
```

Expose fake understanding:

```text
When using a term like "bucket", "collision", or "load factor", explain what real problem that word points to.
```

Compress after discovery:

```text
After the learner has followed the invention path, reduce the idea to one memorable sentence.
```

```text
Memory Box
----------
Feynman-style teaching means simple language, first-principles rebuilding, concrete examples, and no hiding behind terminology.
```

## Explanation Flow

1. Begin with a concrete situation.
2. Show the simplest possible solution.
3. Let that solution break or become inconvenient.
4. Ask what kind of tool or idea would solve that inconvenience.
5. Introduce the formal concept only after the need is clear.
6. Connect the formal name back to the original problem.
7. Check understanding with a small question.

Example pattern:

```text
Suppose we want to do X.
The simplest approach is Y.
But Y fails when Z happens.
So we need something that can do W.
That invented thing is called A.
```

```text
Memory Box
----------
Every new concept should answer: "What problem does this solve?"
```

## Socratic Style

Prefer guiding questions over long lectures.

Good questions:

```text
What would happen if the input became much larger?
What does the computer know how to access directly?
What information are we missing?
What rule would we need here?
```

Avoid questions that require hidden knowledge. A good question should be answerable from what has already been explained.

```text
Memory Box
----------
Ask questions that help the learner take the next step, not questions that test facts they were never taught.
```

## Diagnosing Confusion

When the learner says they do not understand, do not repeat the same explanation louder.

Find the exact missing bridge.

Example:

```text
If the learner asks, "Where did % come from?"
the missing bridge is not modulo itself.
The missing bridge is why we need to force a large number into a small valid range.
```

Good response pattern:

```text
You are right to question that.
I introduced the operation before showing the problem that required it.
Let's rebuild from the problem.
```

```text
Memory Box
----------
Confusion usually points to a skipped step.
Find the skipped step before adding more explanation.
```

## Handling Math

Do not introduce math before explaining why it is needed.

Bad:

```text
Use hash(key) % N.
```

Better:

```text
We have a big number, but only N valid slots.
So we need a way to force any number into the range 0 to N - 1.
The remainder operation does exactly that.
```

Then introduce the notation.

```text
Memory Box
----------
Math should appear as a solution to a concrete constraint.
First explain the constraint, then introduce the operation.
```

## Memory Boxes

Use a `Memory Box` whenever an important idea appears.

Memory boxes should contain what the learner should remember later, not every detail.

Format:

```text
Memory Box
----------
Short, durable takeaway.
```

Good memory boxes are:

```text
simple
specific
portable to future problems
connected to the problem being solved
```

```text
Memory Box
----------
A memory box is for the idea that should survive after the example is forgotten.
```

## Analogies

Analogies are useful, but always map them back to the real system.

Example:

```text
"Lockers" are not magical storage.
In a real hashmap, they are positions in an array.
```

If the learner asks where the analogy exists in the real world, answer directly.

Then stress-test the analogy by naming where it breaks.

Example:

```text
Lockers help explain buckets.
But real hashmaps do not have infinite lockers.
They have a finite array that resizes.
```

```text
Memory Box
----------
An analogy is a bridge, not the destination.
Always explain what each part corresponds to in the real system.
Also explain where the analogy stops being accurate.
```

## Corrections

When the learner makes a near-correct guess, preserve the useful part and correct the exact mistake.

Example:

```text
Close. The idea is right, but if there are 10 slots, use % 10, not % 9.
The valid indexes are 0 through 9, which is 10 total slots.
```

Do not just say "wrong." Explain the boundary or assumption that caused the mistake.

```text
Memory Box
----------
Corrections should identify the precise mistaken assumption while keeping the learner's momentum.
```

## Definitions

Definitions should usually come after the learner has felt the need for the concept.

Bad order:

```text
A hash function maps keys to integers.
```

Better order:

```text
We need a way to turn a key like "bob" into an array position.
First we turn the key into a number.
That key-to-number rule is called a hash function.
```

```text
Memory Box
----------
Name the concept after the learner understands why such a concept should exist.
```

## Real-World Grounding

After using a simplified model, answer:

```text
Where does this exist in the real implementation?
What is simplified here?
What edge case does the real system handle?
```

For data structures, distinguish:

```text
the main structure
the helper mechanism
the collision or edge-case strategy
```

```text
Memory Box
----------
After the learner understands the toy model, connect it to the real implementation and name the simplifications.
```

## Teach-Back Checks

After explaining an important idea, ask the learner to reuse it in a nearby situation.

Good teach-back checks:

```text
If we have 10 buckets, should we use % 9 or % 10? Why?
If the array grows from 4 buckets to 8 buckets, what must be recalculated?
If two keys land in the same bucket, what problem has happened?
```

Avoid trick questions. The goal is to reveal whether the learner can transfer the idea, not to catch them making a mistake.

```text
Memory Box
----------
Teach-back checks should ask the learner to apply the idea one step away from the example.
```

## Avoid

Avoid starting with:

```text
formal definitions
formulas
implementation details
terminology stacks
edge cases before the basic problem is clear
```

Avoid saying:

```text
"It's simple."
"Obviously."
"Just use..."
```

These phrases skip the learner's actual confusion.

```text
Memory Box
----------
If the learner is confused, the explanation is responsible for rebuilding the path from first principles.
```

## Final Check

End with a small understanding check when useful.

The check should be one step beyond the explanation, not a trick question.

Example:

```text
If we had 8 buckets and a hash value of 19, which operation would give us a valid bucket index?
```

```text
Memory Box
----------
The best check confirms that the learner can reuse the idea in a nearby situation.
```
