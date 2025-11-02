#!/usr/bin/env python3
import os
import  subprocess
import  sys
import requests
import  json

def get_staged_diff():
    try:
        result = subprocess.run(
            ['git','diff','--cached','--no-color'],
            capture_output=True,
            text=True,
            check=True
        )
        return  result.stdout
    except subprocess.CalledProcessError as e:
        print(f"Error in getting git diff {e}")
        return ""

def generate_commit_message(diff):
    """Generate commit message using OpenRouter"""
    api_key = os.getenv('OPENROUTER_API_KEY')
    if not api_key:
        print("OPENROUTER_API_KEY not set")
        return None

    if not diff or len(diff.strip()) < 10:
        print("No significant changes detected")
        return None

    limited_diff = diff[:3000]

    prompt = f"""Analyze these git changes and generate a concise, conventional commit message.
Format: type(scope): description

Common types:
- feat: new feature
- fix: bug fix
- docs: documentation
- style: formatting, missing semi colons, etc
- refactor: code restructuring
- test: adding tests
- chore: maintenance tasks

Changes:
{limited_diff}

Return ONLY the commit message in conventional format, nothing else."""

    try:
        response = requests.post(
            "https://openrouter.ai/api/v1/chat/completions",
            headers={
                "Content-Type": "application/json",
                "Authorization": f"Bearer {api_key}",
                "HTTP-Referer": "https://github.com/pre-commit",
                "X-Title": "Pre-commit AI"
            },
            json={
                "model": "anthropic/claude-3.5-sonnet",
                "messages": [
                    {
                        "role": "system",
                        "content": "You are an expert developer that generates perfect conventional commit messages. Be concise and specific."
                    },
                    {
                        "role": "user",
                        "content": prompt
                    }
                ],
                "max_tokens": 100,
                "temperature": 0.1
            },
            timeout=30
        )

        if response.status_code == 200:
            message = response.json()['choices'][0]['message']['content'].strip()
            message = message.split('\n')[0].strip('"\'')
            return message
        else:
            print(f"API error: {response.status_code} - {response.text}")
            return None

    except Exception as e:
        print(f"Error calling OpenRouter: {e}")
        return None

def main():
    print("🚀 AI commit script started")

    # Handle pre-commit argument format: --commit-msg-file <filename>
    if len(sys.argv) > 2 and sys.argv[1] == "--commit-msg-file":
        commit_msg_file = sys.argv[2]
    elif len(sys.argv) > 1:
        # Direct file path (for manual testing)
        commit_msg_file = sys.argv[1]
    else:
        print("❌ No commit message file provided")
        print(f"Arguments received: {sys.argv}")
        return 1

    print(f"📁 Commit message file: {commit_msg_file}")

    if not os.path.exists(commit_msg_file):
        print(f"❌ File does not exist: {commit_msg_file}")
        return 1

    # Read existing commit message
    with open(commit_msg_file, 'r') as f:
        existing_content = f.read().strip()

    print(f"📄 Existing content: {repr(existing_content)}")

    lines = existing_content.split('\n')
    meaningful_lines = [line for line in lines if line.strip() and not line.strip().startswith('#')]

    if meaningful_lines:
        print("ℹ️ Commit message already exists, skipping AI generation")
        return 0

    print("🤖 Analyzing changes with AI...")

    diff = get_staged_diff()
    if not diff:
        print("❌ No staged changes found")
        return 0

    print(f"📊 Diff length: {len(diff)} characters")

    ai_message = generate_commit_message(diff)

    if ai_message:
        with open(commit_msg_file, 'w') as f:
            f.write(ai_message)
        print(f"✅ AI-generated commit message: {ai_message}")
        return 0
    else:
        print("❌ Failed to generate AI commit message")
        return 1

if __name__ == "__main__":
    sys.exit(main())