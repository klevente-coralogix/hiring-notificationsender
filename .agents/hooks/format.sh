#!/usr/bin/env bash
# Formats the Rust and Go files an agent just edited. Used as a post-edit hook
# by Claude Code (.claude/settings.json) and Codex (.codex/hooks.json). Never
# blocks the agent: if a file doesn't parse yet or a tool is missing, it's left
# as is.

input=$(cat)

# Claude Code passes the edited file in tool_input.file_path. Codex passes an
# apply_patch patch in tool_input.command, which names each file in a header.
files=$(jq -r '
  .tool_input.file_path //
  (.tool_input.command // "" | strings
    | split("\n")[]
    | capture("^\\*\\*\\* (Add File|Update File|Move to): (?<path>.+)$").path)
  // empty' <<<"$input" 2>/dev/null) || exit 0
cwd=$(jq -r '.cwd // empty' <<<"$input" 2>/dev/null)

while IFS= read -r file; do
  [ -n "$file" ] || continue
  [[ "$file" = /* ]] || file="${cwd:-.}/$file"
  case "$file" in
    *.rs) rustfmt --edition 2024 "$file" ;;
    *.go) gofmt -w "$file" ;;
  esac >/dev/null 2>&1
done <<<"$files"

exit 0
