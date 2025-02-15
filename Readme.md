# Git Diff to Commit Message Generator

This Go program runs `git diff` to collect changes in a Git repository and then sends the results to the OpenAI API to generate a well-written commit message. It provides an AI-generated commit message that encapsulates the changes made, following best practices for commit messages.

## Features
- Automates commit message generation using OpenAI's API.
- Reads Git repository changes with `git diff`.
- Integrates additional context to create informative and comprehensive commit messages.
- Offers flexibility in choosing the OpenAI model and custom context.
- Customise the format and style of the prompt using `~/.git_diff_prompt.txt`

## Requirements
- Go 1.16 or later
- An OpenAI API key

## Installation

The easiest way is to download a binary for your OS/architecture from the Releases page.  If you want to build the code yourself then :

1. Clone the repository to your local machine:
   ```bash
   git clone https://github.com/ohnotnow/diff_message_generator_go
   ```
2. Change into the project directory:
   ```bash
   cd diff_message_generator_go
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
./gdm --model <MODEL-NAME> [Optionally add some additional context]
```
- **Model**: Specifies the OpenAI model to use (default: `gpt-4o-mini`).
- Optional additional context to be added to the prompt sent to OpenAI.

If you do not specify a model, it defaults to `gpt-4o-mini`. Ensure that you have a valid OpenAI API key set in the environment variable `OPENAI_API_KEY`.

## Output
The output is a well-written commit message based on the `git diff` result, displayed in the console. You can then use this message to commit your changes in Git.
The output is geared to try and follow the 'Conventional Commit' standard.  See https://www.conventionalcommits.org/en/v1.0.0/ .

## Helper shell function

I use this function (MacOS specific) in my `~/.bashrc` to run the gdm command, show the output and also copy it to the clipboard.  

```bash
gdmp() {
    local result

    if [ $# -eq 0 ]; then
        # No arguments, run gdm and pipe to pbcopy
        gdm | pbcopy
        result=$(pbpaste)
    else
        # Arguments provided, capture gdm output then pipe to pbcopy
        result=$(gdm "$@")
        echo "$result" | pbcopy
    fi

    # Print the result to the terminal
    echo "$result"
    echo "The commit message has been copied to your clipboard."
}
```

Then you can just run `gdmp` and do your `git add` and paste in the message.  Or `gdmp improve error handling of missing files` to give some extra context.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing
Contributions are welcome! Please fork this repository and create a pull request with your proposed changes.
