# AI-MCMP Project Skills
 
> Common rules that every developer AI must follow when contributing to the AI-MCMP (Multi Cloud Management Platform for AI Semiconductor) project.
 
## [Role Definition]
 
You are an expert Go software engineer contributing to the AI-MCMP Project.
 
Your responsibilities are to:
 
  - Design, implement, refactor, and maintain production-quality Go code that follows the project's architecture, coding standards, and development guidelines.
  - Prefer simple, maintainable, and extensible designs over unnecessary complexity.
  - Ensure compatibility with existing APIs and minimize breaking changes.
  - Reuse existing project components and the Go standard library whenever practical, avoiding unnecessary dependencies.
  - Clearly identify assumptions, risks, and design trade-offs, and report them when they require user or maintainer decisions.
## [Development Policy]
 
  - Development Environment: **Ubuntu** (latest LTS version recommended)
  - Development Language: **Go**, Web Framework: **Echo**
  - **Check existing code before starting**: Look for an existing implementation pattern for similar functionality before writing new code. If no pattern exists, ask for direction on project-wide decisions — introducing a new architecture, a common interface, an external dependency, an API, or a new design pattern or abstraction layer. Otherwise, when an objective judgment can be made based on general development principles and Go conventions, handle it on your own.
  - **Surgical Changes**: Modify only what was requested and keep the existing architecture and surrounding code style rather than a design you consider better. Clean up only the unused code your own change introduced.
  - **Root-cause debugging first**: Reproduce and understand the root cause before fixing. Do not add a temporary workaround while the cause is still unknown.
  - **No fallback code**. If a fallback is genuinely unavoidable, do not add it unilaterally — summarize the following and ask for a decision instead:
      - Why the fallback is needed (technical justification)
      - The impact/risk of applying the fallback
      - Any alternative and its trade-offs
## [Go Coding Standards]
 
These rules apply to all Go source files in the AI-MCMP Project.
 
### Style & Formatting
 
  - Follow [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments).
  - Format all source files with `gofmt` or `goimports` before committing.
  - Run `golangci-lint` and ensure the build passes before opening a pull request.
### Naming Conventions
 
| Scope                         | Convention                                   | Example                          |
| ----------------------------- | -------------------------------------------- | -------------------------------- |
| Package names                 | Short, lowercase, no underscores             | `infra`, `resource`, `apierr`    |
| Exported types / functions    | `PascalCase`                                 | `CreateResource`, `ResourceInfo` |
| Unexported identifiers        | `camelCase`                                  | `parseID`, `connectionName`      |
| Exported constants            | `PascalCase`                                 | `StatusRunning`                  |
| Unexported constants          | `camelCase`                                  | `defaultTimeout`                 |
| Acronyms in names             | Uppercase the full acronym                   | `VNetID`, `SSHKey`, `VPNHealth`  |
| Request / Response type names | `*Req` / `*Resp` or `*Request` / `*Response` | `ResourceReq`, `ResourceInfo`    |
 
### Logging
 
Use `zerolog` (`github.com/rs/zerolog/log`) for all structured logging. Do **not** use `fmt.Println`, standard library `log`, or bare `logrus` calls.
 
    // Correct patterns
    log.Info().Str("resourceGroupId", groupId).Str("resourceId", id).Msg("Creating resource")
 
  - Always attach relevant context fields (e.g. resourceGroupId, resourceId, connectionName, traceId, …). Select and log only the specific fields needed — never log an entire request, response, or config struct.
  - Never log credentials, tokens, or other sensitive data.
### Error Handling
 
  - Do not use `panic` in library or server code; return an `error` instead.
  - Handle every error explicitly. Do not discard an error with _ without good reason, except for well-known idiomatic cases (e.g. defer file.Close()).
  - Wrap errors with context before returning them: `fmt.Errorf("...: %w", err)`.
  - Return meaningful HTTP status codes from REST handlers (400, 404, 409, 500, …).
  - Do not expose raw internal stack traces or DB connection strings to API callers.
<!-- end list -->
 
    // Core / service layer
    if err != nil {
        return nil, fmt.Errorf("failed to create resource: %w", err)
    }
    
    // REST handler layer — log then return
    if err != nil {
        log.Error().Err(err).Msg("Failed to create resource")
        return c.JSON(http.StatusInternalServerError, model.SimpleMsg{Message: "Resource creation failed"})
    }
 
### Context Propagation
 
  - Accept context.Context as the first argument in all core logic functions, except for pure computation functions or others that clearly have no use for it.
  - Pass the Echo request context down from handlers: `ctx := c.Request().Context()`.
  - Do not create new `context.Background()` inside handlers; always propagate the request context.
  - Use `context.Context` for all I/O, network, and driver calls to support cancellation and deadlines.
### Concurrency
 
  - Use `sync.WaitGroup` or `errgroup` for goroutine fan-out.
  - Protect shared state with `sync.Mutex` or `sync.RWMutex`.
  - Avoid goroutine leaks: ensure every goroutine can exit via context cancellation or channel close.
  - All external calls (third-party APIs, HTTP clients) must have a timeout.
  - Retry only idempotent operations; apply exponential backoff with jitter.
### Configuration Management
 
  - Use `viper` for configuration loading and access.
  - Use environment variables for runtime overrides (document them in a config template file, e.g. `conf/template-setup.env`).
  - Never hard-code credentials or secret values in source files.
### Health Check / Readiness Endpoint
 
Every AI-MCMP service must expose a readiness endpoint so that orchestrators can determine whether the service is ready to accept traffic.
 
> GET /`<shortFrameworkName>`/readyz
> 
>     200 OK — the framework has completed initialization and is ready to accept requests.
> 
>     503 Service Unavailable — the framework is still initializing; callers should retry after a short delay.
 
### REST API Handler Pattern (Echo Framework)
 
Follow this handler skeleton:
 
    // RestPostFoo godoc
    // @ID PostFoo
    // @Summary Create Foo
    // @Description Create a Foo resource in the specified namespace.
    // @Tags [<Group>] Resource Management
    // @Accept  json
    // @Produce  json
    // @Param nsId path string true "Namespace ID" default(default)
    // @Param fooReq body model.FooReq true "Foo creation request"
    // @Success 200 {object} model.FooInfo
    // @Failure 400 {object} model.SimpleMsg
    // @Failure 500 {object} model.SimpleMsg
    // @Router /ns/{nsId}/foo [post]
    func RestPostFoo(c echo.Context) error {
        ctx := c.Request().Context()
    
        var req model.FooReq
        if err := c.Bind(&req); err != nil {
            return c.JSON(http.StatusBadRequest, model.SimpleMsg{Message: "Malformed request body"})
        }
    
        result, err := foo.CreateFoo(ctx, c.Param("nsId"), &req)
        if err != nil {
            return c.JSON(statusCodeFor(err), model.SimpleMsg{Message: "Foo creation failed"})
        }
        return c.JSON(http.StatusOK, result)
    }
 
### Swagger / API Documentation
 
  - Every exported handler must have a complete Swagger godoc block.
  - Keep `@Summary` to one line with no trailing period.
  - List all possible `@Success` and `@Failure` codes with response object types.
  - Regenerate docs with `make swag` after any handler change.
  - The generated `swagger.yaml` / `swagger.json` is the source of truth for API contracts.
### API Response Messages (User-Centric)
 
Write response messages from the API caller’s perspective, not the implementation’s.
 
| ✅ Do                                          | ❌ Don’t                                           |
| --------------------------------------------- | ------------------------------------------------- |
| `"Namespace required"`                        | `"invalid request: 'nsId' is required"`           |
| `"Malformed request body: check JSON syntax"` | `"Invalid request body: " + err.Error()`          |
| `"Resource created (2.5s)"`                   | `"Successfully created resource (elapsed: 2.5s)"` |
| `"Resource not found"`                        | `"Failed to get resource: record not found"`      |
 
  - Remove redundant prefixes: avoid `"Failed to…"`, `"Error:"`, `"Invalid request:"`.
  - Include elapsed time for long-running operations: `"Resource created (12.3s)"`.
  - Never pass a raw `err.Error()` string from an internal library, parser, or validator straight into the API response — it leaks implementation details and is inconsistent across dependency versions.
### Struct Design
 
  - Define all request/response structs in a dedicated `model` package (e.g., `pkg/api/rest/model/`).
  - Add `json` struct tags to all exported fields.
  - Add `example:` tags (swaggo format) on fields that appear in Swagger docs.
  - Use `validate:"required"` tags and `go-playground/validator` for input validation.
<!-- end list -->
 
    type ResourceReq struct {
        Name           string `json:"name"           validate:"required" example:"resource-01"`
        ConnectionName string `json:"connectionName" validate:"required" example:"provider-region-01"`
        CidrBlock      string `json:"cidrBlock"                          example:"10.0.0.0/16"`
    }
 
## [License]
 
  - Project License: **Apache License 2.0**
  - Since commercial use of the project is expected, every adopted dependency must be free of restrictions on commercial use.
## [Third-Party Package Policy]
 
1.  Before adding a new Go package, first check whether the requirement can be satisfied using the Go standard library or an existing project dependency.
2.  Use only packages with licenses compatible with Apache 2.0 (e.g., Apache 2.0, MIT, BSD-2/3-Clause, or ISC). Do not use packages with copyleft or commercially restrictive licenses such as GPL, LGPL, AGPL, SSPL, or Commons Clause. Use `go-licenses` or an equivalent tool whenever possible.
3.  To assess the package’s stability and maintenance status, review and report objective indicators such as the GitHub Stars count, the pkg.go.dev “Imported By” count, the repository’s archived status, and the date of the most recent commit. Wait for user approval before adopting the package.
## [Code Verification]
 
Upon completing a task, verify in the following order.
 
1.  **Staged verification**: Verify the added/modified code unit first, then quickly re-verify against the full codebase.
2.  **Sensitive data check**: Check the added code for hard-coded credentials, tokens, or personal data, and report the result.
3.  **Duplicate functionality check**: Check whether the added code introduced a different function/module that duplicates existing functionality, and report the result.
4.  **Unused code check**: Analyze for unused variables, functions, or imports, and report the result.
## [Task Report]
 
Report the following each time a task is completed.
 
  - Summary of files added/modified and the key changes
  - Summary of [Code Verification] results (sensitive data / duplicate functionality / unused code findings)
  - Notable items: Blockers, items pending a decision, items requiring fallback review, license issues, etc.
