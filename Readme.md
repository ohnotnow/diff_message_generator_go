# Git Diff to Commit Message Generator

This Go program runs `git diff` to collect changes in a Git repository and then sends the results to the OpenAI API to generate a well-written commit message. It provides an AI-generated commit message that encapsulates the changes made, following best practices for commit messages.

## Features
- Automates commit message generation using OpenAI's API.
- Reads Git repository changes with `git diff`.
- Integrates additional context to create informative and comprehensive commit messages.
- Offers flexibility in choosing the OpenAI model and custom context.

## Requirements
- Go 1.16 or later
- An OpenAI API key

## Installation
1. Clone the repository to your local machine:
   ```bash
   git clone <YOUR-REPO-URL>
   ```
2. Change into the project directory:
   ```bash
   cd <YOUR-REPO-NAME>
   ```
3. Build the Go program:
   ```bash
   go build gdm.go
   ```

## Configuration
- **OpenAI API Key**: Set your OpenAI API key as an environment variable:
  ```bash
  export OPENAI_API_KEY=<YOUR-API-KEY>
  ```

## Usage
To run the program, use the following command:
```bash
./gdm --model <MODEL-NAME> --context "Your additional context"
```
- **Model**: Specifies the OpenAI model to use (default: `gpt-4-turbo`).
- **Context**: Optional additional context to be added to the prompt sent to OpenAI.

If you do not specify a model, it defaults to `gpt-4-turbo`. Ensure that you have a valid OpenAI API key set in the environment variable `OPENAI_API_KEY`.

## Output
The output is a well-written commit message based on the `git diff` result, displayed in the console. You can then use this message to commit your changes in Git.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing
Contributions are welcome! Please fork this repository and create a pull request with your proposed changes.
