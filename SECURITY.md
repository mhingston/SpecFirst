# Security Policy

SpecFirst is a tool that processes your specifications and code. Because it often interfaces with Large Language Models (LLMs), it is critical to handle sensitive data with care.

## Data Handling & Privacy

SpecFirst itself does not send data anywhere unless you configure it to pipe output to an external tool (like an LLM CLI). However, the workflows you build often involve sending code and prompts to third-party services.

> [!WARNING]
> **Do not include secrets, credentials, PII (Personally Identifiable Information), or sensitive customer data in your prompts.**

Always review the context files you include in your `inputs:` section.

## Safe Scoping Defaults

When defining the scope of a task (e.g., in a `scope.md` template or `protocol.yaml`), explicitly exclude sensitive files.

**Recommended Exclusion List:**
- `.env`, `.env.*`
- `*.pem`, `*.key`, `*.p12`, `*.kdbx`
- `id_rsa*`, `id_dsa*`
- `credentials*.json`
- `secrets*`
- `*_token.txt`

## Third-Party Model Caution

If you pipe SpecFirst output to hosted LLMs (e.g., Gemini, OpenAI, Claude):
1.  **Understand Retention**: Check if the provider retains your data for training.
2.  **Compliance**: Ensure you are compliant with your organization's data policies (GDPR, SOC2, etc.).
3.  **Opt-Out**: Use "Enterprise" or "Zero Data Retention" tiers where available.

## Local/Offline Support

SpecFirst is fully functional offline.
- It runs locally on your machine.
- It does not require an internet connection to generate prompts or render templates.
- If you use a local LLM runner (like Ollama or LocalAI), the entire loop can remain air-gapped.
