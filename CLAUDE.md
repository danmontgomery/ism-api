# ISM API — Development Rules

## TDD Workflow

Start with tests. Follow red/green TDD:
1. Write a failing test first
2. Write the minimum code to make it pass
3. Refactor if needed
4. Repeat

Do not write implementation code without a corresponding failing test already in place.

## Environment Variables

Document every environment variable the application reads in README.md. Each entry must include:
- Variable name
- Description
- Default value (if any)
- Whether it is required or optional
