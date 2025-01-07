# Autolang GitHub Action

## 📖 **Overview**

This GitHub Action, **Autolang**, automates the translation of files in your repository using OpenAI's GPT API. It scans specified folders, detects files written in a base language, and generates translations for each target language specified. The action is fully customizable and integrates seamlessly with your CI/CD pipeline.

---

## 🚀 **Features**

- Supports multiple source folders.
- Translates content from a base language to multiple target languages.
- Generates translated files with proper naming conventions.
- Removes old translated files and regenerates new ones to ensure up-to-date translations.
- Cleans unnecessary block delimiters (```), making it compatible with .yml, .html, .json, .xml, and other formats.
- Powered by OpenAI's GPT models for accurate and context-aware translations.
- Automates the creation of multilingual content during deployments.

---

## 🧩 **Inputs**

| Input Name            | Description                                 | Required | Default |
| --------------------- | ------------------------------------------- | -------- | ------- |
| `translations_folder` | Comma-separated list of folders to process. | ✅        |         |
| `base_language`       | The language of the original files.         | ✅        |         |
| `target_languages`    | Comma-separated list of target languages.   | ✅        |         |
| `openai_api_key`      | Your OpenAI API key.                        | ✅        |         |

### Example:

```yaml
inputs:
  translations_folder: "src/translations,docs"
  base_language: "en"
  target_languages: "es,fr,de"
  openai_api_key: ${{ secrets.OPENAI_API_KEY }}
```

---

## ⚙️ **Environment Variables**

Make sure you have the following environment variables set:

| Variable Name    | Description          |
| ---------------- | -------------------- |
| `OPENAI_API_KEY` | Your OpenAI API key. |

---

## 🛠️ **Usage**

### **Basic Example:**

```yaml
name: Translate Files

on:
  push:
    branches:
      - main

jobs:
  translate:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v3

    - name: Run Autolang Action
      uses: your-repo/autolang-action@v1
      with:
        translations_folder: "src/translations,docs"
        base_language: "en"
        target_languages: "es,fr,de"
      env:
        OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
```

### **Advanced Example:**

```yaml
name: Multilingual Deployment

on:
  workflow_dispatch:

jobs:
  translate:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v3

    - name: Install dependencies
      run: go mod download

    - name: Run Autolang
      uses: your-repo/autolang-action@v1
      with:
        translations_folder: "src/translations,content"
        base_language: "en"
        target_languages: "es,fr,ja,zh"
      env:
        OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
```

---

## 📁 **File Naming Convention**

The action will create translated files following this pattern:

```
<OriginalFileName>.<LanguageCode><Extension>
```

For example:

- `index.html` -> `index.es.html`
- `content.md` -> `content.fr.md`

### Important Note:

The action ensures that translated files do not contain multiple language codes in their names. For instance:

- `EN.ES.html` or `example.EN.ES.html` will be corrected to `ES.html` or `example.ES.html`.

---

## 🧪 **How It Works**

1. The action reads the `translations_folder` input and processes each specified folder.
2. It scans the files in each folder and reads their content.
3. For each target language, it sends a request to the OpenAI API to translate the file content.
4. The translated content is cleaned to remove any unnecessary block delimiters (```).
5. The translated content is saved as a new file in the same directory, with a language-specific suffix added to the filename.
6. The action removes old translated files and regenerates them to ensure they are always up-to-date.
7. The action repeats the process for all specified folders and languages.

---

## 🧵 **Process Flow**

1. **Input Validation**: Ensures that the required inputs are provided.
2. **Folder Scanning**: Walks through each specified folder to find files.
3. **File Translation**: Sends file content to the OpenAI API for translation.
4. **File Cleaning**: Removes unnecessary block delimiters (```).
5. **File Creation**: Saves the translated content as new files in the same directory.
6. **Completion**: Logs a success message once all translations are complete.

---

## 🔍 **Error Handling**

- If any input is missing, the action will terminate and log an appropriate error message.
- If a file cannot be read or written, an error message will be logged with the file path.
- If the OpenAI API returns an error, it will be caught and logged.

---

## 🛡️ **Security Considerations**

- Store your `OPENAI_API_KEY` as a GitHub Secret to avoid exposing it in your workflow file.
- Ensure your API key has the necessary permissions to access the OpenAI API.

---

## 📚 **Dependencies**

- [go-openai](https://github.com/sashabaranov/go-openai): A Go client for accessing OpenAI's API.

---

## 📄 **License**

This project is licensed under the MIT License. See the `LICENSE` file for more details.

---

## 👨‍💻 **Contributing**

Contributions are welcome! Please follow the standard GitHub flow for pull requests:

1. Fork the repository.
2. Create a new branch.
3. Make your changes.
4. Submit a pull request.

---

## 🤝 **Support**

If you encounter any issues or have any questions, feel free to open an issue in the GitHub repository.

---

## 📬 **Contact**

- **Author**: Edison J. Padilla
- **Email**: edisonpadilla.dev@gmail.com

