# AI-Commits: The Commit Message Generator 

> 🚀 Write meaningful commit messages in seconds with AI-powered suggestions!

AI-Commits is a command-line tool that generates conventional commit messages by analyzing your staged changes. It understands your code changes and explanation, then generates a clear, descriptive commit message in English following the Conventional Commits format.

[Demo Video Coming Soon]

##  Installation
### Prerequisites

- Go 1.20 or later must be installed. [Install Go](https://go.dev/doc/install)

### ✅ Option 1: Install via `go install` (recommended)

```bash
go install github.com/canermastan/ai-commits@latest
```

After installation, make sure $GOPATH/bin (usually ~/go/bin) is in your PATH, so you can use it globally:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then simply run:

```bash
ai-commits
```

### ⚡ Option 2: Build manually

```bash
# Clone the repository
git clone https://github.com/canermastan/ai-commits.git

# Navigate to the project directory
cd ai-commits

# Build the project
go build

# Optional: Add to your PATH for global usage
```

## 🔑 Configuration

1. Get your Gemini API key from [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Set the environment variable:

```bash
# For Windows PowerShell
 $env:GEMINI_API_KEY = "your-api-key"

# For Linux/macOS
export GEMINI_API_KEY="your-api-key"
```

## 💻 Usage

```bash
# Normal mode (interactive)
ai-commits

# Fast mode (non-interactive)
ai-commits --fast
```

### Example Flow
1. Stage your changes with `git add`
2. Run `ai-commits`
3. Explain what you did
4. Review and confirm the generated commit message
5. Done! 🎉

## 🛣️ Roadmap

- [ ] Enhanced UI with more interactive elements
- [ ] Support for Local LLM integration
- [ ] Custom commit message templates
- [ ] Batch commit message generation
- [ ] Support for more AI providers
- [ ] Commit message history

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.