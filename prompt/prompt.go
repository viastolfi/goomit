package prompt

import (
	"os"
	"strings"

	"main/client"
)

func getPromptSpecificContext() string {
	data, err := os.ReadFile(".goomit/context.md")
	if err != nil {
		return "" // context is optional
	}
	return string(data)
}

func commitPromptTemplate(projectContext, diff string) string {
	var b strings.Builder
	b.WriteString("ROLE: You are an experienced software engineer. Craft a high-quality Git commit message.\n")
	if projectContext != "" {
		b.WriteString("PROJECT CONTEXT (for understanding only, do NOT copy verbatim):\n" + projectContext + "\n")
	}
	b.WriteString(`INSTRUCTIONS:
1. Use Conventional Commits. Allowed types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert.
2. Header format: <type>(<optional-scope>): <imperative short summary>
   - Max 72 chars (hard limit). No trailing period. Imperative mood ("add", "fix", "refactor").
   - Pick ONE best type (if only dependency / tooling changes => chore; only test code => test; performance improvement => perf; code re-organisation without behaviour change => refactor).
3. If multiple conceptual changes that should be separate commits are present, STILL produce ONE message capturing the dominant theme (prefer the user curates their diff). Do NOT invent multiple commits.
4. Body (optional) only if it adds value beyond the header. When used:
   - Start after a blank line.
   - Provide concise bullet points starting with '- ' each.
   - Each bullet: high-level intent / effect / rationale. DO NOT paraphrase code line-by-line.
   - Group related file changes together; mention file paths only when necessary for clarity (prefer patterns, e.g. 'client/' or 'prompt:').
5. MUST NOT include: code blocks, markdown headers, backticks, numbered lists, emojis, hashtags, quotes, HTML, or extra commentary outside the commit message.
6. Detect special cases:
   - Pure revert (diff shows previous changes removed) => use revert: and header "revert: <original header>" if recoverable.
   - Added tests only => test:
   - Performance-focused change (algorithmic / reduced complexity) => perf:
   - Only formatting / whitespace => style:
7. Security / bug fix: briefly indicate root cause or surface impact (e.g. 'fix(api): prevent nil pointer on empty diff').
8. Do NOT fabricate information not evident from diff/context.
9. Ignore ANSI color codes or diff metadata noise if any.
10. If diff is extremely large (> ~2,000 changed lines) summarise major themes (e.g. 'migrate X', 'rename Y') instead of enumerating every file.

OUTPUT REQUIREMENTS:
- Output ONLY the commit message text (header + optional body). Nothing else.
- No leading or trailing blank lines beyond the standard single blank line separating header and body.

REFERENCE EXAMPLES (do NOT copy literally):
feat(auth): add OAuth2 login flow

- introduce token exchange with provider X
- persist refresh token securely in encrypted store
- update user model with provider id

fix(cli): handle empty config path

- prevent nil dereference when no config file exists
- add user-facing error with hint to run 'init'

refactor(prompt): consolidate template building logic

- merge scattered string concatenations into builder pattern
- no functional changes

Now analyse the provided diff and produce the commit message.
`)
	b.WriteString("\nDIFF START\n" + diff + "\nDIFF END\n")
	b.WriteString("\nRETURN ONLY THE COMMIT MESSAGE (header + optional body).")
	return b.String()
}

func promptConfContext() string {
	return "Your task is to summarise multiple sources about a Git repository so the summary can later feed another LLM when generating commit messages. You will be given README content plus basic repo metadata (language, description). Produce concise markdown with:\n* Very short app summary (1 sentence)\n* Key features (bullets)\n* Brief tech stack (bullets, only essentials)\nKeep it short, factual, no marketing fluff.\n"
}

func GeneratePrompt() (string, error) {
	projectContext := getPromptSpecificContext()
	diff, err := client.GetGitDiff()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(commitPromptTemplate(projectContext, diff)), nil
}

func GenerateConfPrompt(context string) string {
	return strings.TrimSpace(promptConfContext() + context)
}
