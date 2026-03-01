---
name: start-task
description: "Start working on a GitHub Issue end-to-end: fetch the issue, plan with task-progress-guide.md, implement via TDD (test-writer, code-implementer, code-reviewer agents), commit, and open a PR."
argument-hint: "<issue-number>"
user-invokable: true
---

# Start Task Skill

## Context

- Current branch: !`git branch --show-current`
- Recent commits: !`git log --oneline -5`
- Open tasks: !`echo "(task list will be populated during execution)"`

## Workflow

### Step 0: Parse arguments

`$ARGUMENTS` must contain a GitHub Issue number (e.g. `42` or `#42`).
Extract the numeric part and assign it to `ISSUE_NUMBER`.
If no number is found, stop and ask the user: "Which issue number should I start on?"

### Step 1: Fetch the Issue

Run the following to get full issue details:

```bash
gh issue view $ISSUE_NUMBER --json number,title,body,labels,assignees
```

Read the output carefully:
- Understand the **summary**, **background**, **goals**, and **acceptance criteria** (完了条件).
- Identify which layers are touched: Backend (Go), Frontend, DB migration, Docs.
- List the concrete sub-tasks from the "やること" section of the issue body.

If the issue body is missing or too vague, stop and ask the user to clarify before continuing.

### Step 2: Create a feature branch

Branch name format: `feature/<ISSUE_NUMBER>-<short-slug>`

Where `<short-slug>` is a lowercase, hyphen-separated English summary of the issue title (max 5 words).

```bash
git checkout -b feature/<ISSUE_NUMBER>-<short-slug>
```

### Step 3: Register tasks in the task list

Use `TaskCreate` to register each sub-task identified in Step 1.
Use `TaskUpdate` to set up dependency chains (`addBlockedBy`) where ordering matters (e.g., domain logic before HTTP handler, backend before frontend).

Always create at minimum:
- One task per layer touched (Backend, Frontend, DB, Docs)
- A final "commit & PR" task that is blocked by all implementation tasks

### Step 4: Implement following TDD (Red → Green → Refactor)

Consult `docs/guidelines/task-progress-guide.md` for the full TDD flow.
Process each implementation task in dependency order.

For each task, execute the following sub-steps **sequentially**:

#### 4a. Test-first (Red)

Invoke the **test-writer** agent:

```
Agent: test-writer
Prompt: "Write failing tests for the following task based on the acceptance
criteria. Task: <task subject and description>. Issue context: <issue body
excerpt>. Follow TDD red phase — tests must fail at this point."
```

Verify tests fail before moving on (run `make test-unit` or the relevant
test command). Mark the task as `in_progress`.

#### 4b. Implement (Green)

Invoke the **code-implementer** agent:

```
Agent: code-implementer
Prompt: "Implement the minimum code to make the failing tests pass.
Task: <task subject and description>. Issue context: <issue body excerpt>.
Follow the repository architecture in go-backend/CLAUDE.md. Tests are
already written — make them green."
```

Run `make test-unit` (and `make test-integration` if integration paths
changed). Confirm all tests pass.

#### 4c. Refactor & review (Refactor)

Invoke the **code-reviewer** agent:

```
Agent: code-reviewer
Prompt: "Review the implementation for this task. Check: correctness,
security (OWASP top 10), performance (N+1, timeouts), readability, adherence
to repository guidelines in docs/guidelines/backend-coding-guidline.md, and
test coverage. Task: <task description>."
```

Apply any critical feedback from the reviewer before marking the task complete.

Run the full quality gate:
```bash
make fmt && make lint && make test-unit
```

If integration paths were touched, also run:
```bash
make test-integration
```

Mark the task as `completed` once all checks pass.

#### 4d. Repeat for each task

Continue until every implementation task is `completed`.

### Step 5: Final quality gate

Run the complete suite one last time and confirm all targets pass:

```bash
make fmt && make lint && make test-unit
```

If any target fails, investigate and fix before proceeding.

### Step 6: Commit

Use the **commit** skill to commit all changes.
Pass the issue number as argument so it appears in every commit message.

```
Skill: commit
Args: $ISSUE_NUMBER
```

### Step 7: Open a Pull Request

Use the **create-pr** skill to open the PR against `main`.

```
Skill: create-pr
```

The PR description must follow the template in the create-pr skill (Japanese body) and include `Closes #<ISSUE_NUMBER>`.

---

## Safety rules

- Never start on `main` — always create a feature branch first.
- Never skip tests. If tests cannot be written (e.g., pure config change), document why.
- Never commit secrets or credentials.
- If any step fails or produces unclear results, stop and report to the user before continuing.
- Use `TaskUpdate` to keep task statuses current throughout the workflow.
