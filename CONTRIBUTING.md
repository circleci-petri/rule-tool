# Contributing to this project

This project is intended to be an experiment in agentic development, where the majority of the code is written and maintained by AI agents rather than human developers. Contributions should focus on guiding and improving the agents' ability to develop the project, rather than directly writing code yourself.

While it is possible to manually edit or refactor the code, this is generally discouraged as it can quickly lead to inconsistencies with the agent-driven development process. Instead, focus on improving the rules and instructions that guide the agents.

## The `rules/` folder

The [`rules/`](rules/) folder contains example rules that ship with this project. These rules are intended to demonstrate how to write rules and can also serve as prompts to guide the agents in specific development tasks. Feel free to explore and modify these rules, or add your own, to experiment with different approaches to agentic development.

# Tips for debugging rules as you work on them

## Visibility

Trying to make sure a rule is always being included?

Ask the LLM to tell you about their required_instructions. You should get some output like

```
The user is asking me about the required_instructions that were provided in the custom_instructions section. Let me analyze what's in there:
In the required_instructions section, there is one rule:
...
```
