# Contributing to fwsync

Thank you for your interest in contributing to fwsync! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful environment for all contributors.

Just don't be a jerk.

## Getting Started

### Prerequisites

- Go 1.24 or higher
- [task](https://github.com/go-task/task)
- [golangci-lint](https://github.com/golangci/golangci-lint)

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/fwsync.git
   cd fwsync
   ```
3. Add the upstream repository as a remote:
   ```bash
   git remote add upstream https://github.com/jharshman/fwsync.git
   ```
4. Build the project:
   ```bash
   go-task all
   ```

## Development Workflow

### Creating a Branch

Create a new branch for your work:
```bash
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/` for new features
- `fix/` for bug fixes
- `docs/` for documentation changes
- `refactor/` for code refactoring

### Making Changes

1. Make your changes in your feature branch
2. Write or update tests as needed
3. Ensure your code follows the project's coding standards
4. Test your changes locally

### Running Tests

Run the test suite to ensure your changes don't break existing functionality:
```bash
go-task test
```

### Linting

This project uses `golangci-lint` for code quality checks. Ensure your code passes linting:

```bash
go-task lint
```

All code must pass linting and tests before being merged.

### Commit Messages

Write clear and meaningful commit messages:
- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Keep the first line under 50 characters
- Add a detailed description if necessary after a blank line
- Include a relevant issue number if applicable

Example:
```
Add support for multiple firewall rules

Extend config to support multiple rules.
Update the init and update commands to support
config changes.

Resolves: ISSUENUMBER
```

## Submitting Changes

### Pull Request Guidelines

- Keep PRs focused on a single feature or fix
- Ensure all tests pass before requesting review
- Update documentation if you're changing functionality
- Add tests for new features or bug fixes
- Squash your commits before merging

## Reporting Issues

### Bug Reports

When reporting bugs, please include:
- A clear, descriptive title
- Steps to reproduce the issue
- Expected behavior vs actual behavior
- Your environment (OS, Go version, Provider configuration)
- Relevant logs or error messages
- Screenshots if applicable

### Feature Requests

For feature requests, please include:
- A clear description of the feature
- Use cases and why it would be valuable
- Any implementation ideas you have
- Examples of how it would work

### Testing

- Write unit tests for new functionality
- Aim for good test coverage
- Mock external dependencies
- Test error cases, not just happy paths

