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
    commit_msg_file = sys.argv[1] if len(sys.argv) > 1 else None

    if not commit_msg_file:
        print("No commit message file provided")
        return 1
    with open(commit_msg_file, 'r') as f:
        existing_content = f.read().strip()
    if existing_content and not existing_content.startswith('#'):
        print("Commit message already exists, skipping AI generation")
        return 0

    print("Analyzing changes with AI...")
    diff = get_staged_diff()
    if not diff:
        print("No staged changes found")
        return 0

    ai_message = generate_commit_message(diff)

    if ai_message:
        with open(commit_msg_file, 'w') as f:
            f.write(ai_message)
        print(f"✅ AI-generated commit message: {ai_message}")
    else:
        print("❌ Failed to generate AI commit message")

    return 0

if __name__ == "__main__":
    sys.exit(main())