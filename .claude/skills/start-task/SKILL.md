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
- Understand the **summary**, **background**, **goals**, and **acceptance criteria** (the "完了条件" section).
- Identify which layers are touched: Backend (Go), Frontend, DB migration, Docs.
- List the concrete sub-tasks from the "やること" (To-do) section of the issue body.

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
- An E2E task if UI or user-facing API endpoints are added or changed (see `.claude/skills/start-task/e2e-context.md`)
- A final "commit & PR" task that is blocked by all implementation tasks

### Step 4: Implement following TDD (Red → Green → Refactor)

Consult `docs/guidlines/task-progress-guide.md` for the full TDD flow.
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

Refer to the context file for the task's layer and invoke the code-implementer agent using the "Implement prompt" section:
- Backend (Go): `.claude/skills/start-task/backend-context.md`
- Frontend (React): `.claude/skills/start-task/frontend-context.md`

Confirm all tests pass before moving on.

#### 4c. Refactor & review (Refactor)

Refer to the context file for the task's layer and invoke the code-reviewer agent using the "Review prompt" section:
- Backend (Go): `.claude/skills/start-task/backend-context.md`
- Frontend (React): `.claude/skills/start-task/frontend-context.md`

Apply any critical feedback, then run the commands in the "Quality gate" section of the same context file.

Mark the task as `completed` once all checks pass.

#### 4e. Process the E2E task (if created in Step 3)

If an E2E task was registered, process it like any other task.
Refer to `.claude/skills/start-task/e2e-context.md` for the test-writer prompt and test execution commands.

#### 4f. Repeat for each task

Continue until every implementation task is `completed`.

### Step 5: Final quality gate

For each layer touched by this PR, refer to its context file and run all commands in the "Quality gate" section:
- Backend (Go): `.claude/skills/start-task/backend-context.md`
- Frontend (React): `.claude/skills/start-task/frontend-context.md`

If any target fails or coverage regresses, investigate and fix before proceeding.

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
