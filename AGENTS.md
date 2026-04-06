# ISM API — Project Rules

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

## Non-Interactive Shell Commands

**ALWAYS use non-interactive flags** with file operations to avoid hanging on confirmation prompts.

Shell commands like `cp`, `mv`, and `rm` may be aliased to include `-i` (interactive) mode on some systems, causing the agent to hang indefinitely waiting for y/n input.

**Use these forms instead:**
```bash
# Force overwrite without prompting
cp -f source dest           # NOT: cp source dest
mv -f source dest           # NOT: mv source dest
rm -f file                  # NOT: rm file

# For recursive operations
rm -rf directory            # NOT: rm -r directory
cp -rf source dest          # NOT: cp -r source dest
```

**Other commands that may prompt:**
- `scp` - use `-o BatchMode=yes` for non-interactive
- `ssh` - use `-o BatchMode=yes` to fail instead of prompting
- `apt-get` - use `-y` flag
- `brew` - use `HOMEBREW_NO_AUTO_UPDATE=1` env var
