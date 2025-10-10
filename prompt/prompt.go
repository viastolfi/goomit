package prompt

import (
	"strings"

	"main/client"
)

// Comming from : https://github.com/GNtousakis/llm-commit/blob/main/llm_commit.py#L181
func promptContext() string {
	return "You are a professional developer with more than 20 years of experience. You're an expert at writing Git commit messages from code diffs. Focus on highlighting the added value of changes (meta-analysis, what could have happened without this change?), followed by bullet points detailing key changes (avoid paraphrasing). Use the specified commit Git style, while forbidding other syntax markers or tags (e.g., markdown, HTML, etc.)"
}

func promptConfContext() string {
	return "Your task is to summarize a bunch of different informtation about a git repository so this summization can then be used as a context element for other llm prompt to generate commit Message. You are going to be gived a README.md file text, the language used in the code base and the git description of the repository. Write your out puts as an md text that will then be writted in a file. Keep it not to verbose and relatively short. Make it so it's usable by llm to get context. I want the following things :\n* Very short summarisation of the apps\n* The key features\n* A brief tech stack with only the essential information\n"
}

func promptTitle() string {
	return "Generate a concise commit message starting with a type keyword (fix:, feat:, ci:, etc.) followed by a one-line summary describing the overall change clearly and briefly. Keep it short, precise, and focused."
}

func promptEnding() string {
	return "just give me the actual commit message and nothing else in your answer. I want to use your response as it is."
}

func GeneratePrompt() (string, error) {
	diff, err := client.GetGitDiff()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(promptContext() + promptTitle() + "Here is the git diff I want you to generate a commit message for : " + diff + promptEnding()), nil
}

func GenerateConfPrompt(context string) string {
	return strings.TrimSpace(promptConfContext() + context)
}
