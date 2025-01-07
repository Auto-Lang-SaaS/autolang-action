package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func main() {
	foldersEnv := os.Getenv("translations_folder")
	if foldersEnv == "" {
		fmt.Println("Error: translations_folder not set")
		os.Exit(1)
	}

	baseLanguage := os.Getenv("base_language")
	if baseLanguage == "" {
		fmt.Println("Error: base_language not set")
		os.Exit(1)
	}

	targetLanguages := os.Getenv("target_languages")
	if targetLanguages == "" {
		fmt.Println("Error: target_languages not set")
		os.Exit(1)
	}
	languages := strings.Split(targetLanguages, ",")

	apiKey := os.Getenv("openai_api_key")
	if apiKey == "" {
		fmt.Println("Error: openai_api_key not set")
		os.Exit(1)
	}
	client := openai.NewClient(apiKey)

	// Dividir rutas proporcionadas en múltiples carpetas
	folders := strings.Split(foldersEnv, ",")

	// Eliminar todos los archivos que no sean del idioma base antes de crear nuevos
	replaceTranslations(folders, baseLanguage, languages, client)

	fmt.Println("Translation completed successfully!")
}

func replaceTranslations(folders []string, baseLanguage string, languages []string, client *openai.Client) {
	// Crear un nuevo estado eliminando traducciones antiguas
	for _, folder := range folders {
		folder = strings.TrimSpace(folder)
		if err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			ext := filepath.Ext(path)
			baseName := strings.TrimSuffix(filepath.Base(path), ext)

			// Mantener solo los archivos del idioma base
			if strings.HasSuffix(baseName, "."+baseLanguage) || strings.EqualFold(baseName+ext, baseLanguage+ext) {
				return nil
			}

			// Eliminar archivos traducidos antiguos
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove old translation file %s: %w", path, err)
			}
			fmt.Printf("Removed file: %s\n", path)
			return nil
		}); err != nil {
			fmt.Printf("Error removing old translations in folder %s: %v\n", folder, err)
		}
	}

	// Crear nuevas traducciones
	for _, folder := range folders {
		folder = strings.TrimSpace(folder)
		if err := processFolder(folder, baseLanguage, languages, client); err != nil {
			fmt.Printf("Error processing folder %s: %v\n", folder, err)
		}
	}
}

func processFolder(folder, baseLanguage string, languages []string, client *openai.Client) error {
	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignorar directorios, procesar solo archivos
		if info.IsDir() {
			return nil
		}

		// Leer el contenido del archivo
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Obtener la extensión y el nombre del archivo
		ext := filepath.Ext(path)
		baseName := strings.TrimSuffix(filepath.Base(path), ext)

		// Eliminar cualquier sufijo de idioma existente del nombre base
		if strings.HasSuffix(baseName, "."+baseLanguage) {
			baseName = strings.TrimSuffix(baseName, "."+baseLanguage)
		}

		// Traducir y crear archivos para cada idioma
		for _, lang := range languages {
			if lang == baseLanguage {
				continue
			}

			translatedContent, err := translateText(string(content), lang, client)
			if err != nil {
				return fmt.Errorf("failed to translate file %s to %s: %w", path, lang, err)
			}

			// Limpiar los delimitadores de bloques (```)
			translatedContent = cleanBlockDelimiters(translatedContent)

			// Manejar correctamente el nombre del archivo traducido
			var newFileName string
			if baseName == baseLanguage {
				newFileName = fmt.Sprintf("%s%s", lang, ext)
			} else {
				newFileName = fmt.Sprintf("%s.%s%s", baseName, lang, ext)
			}

			newFilePath := filepath.Join(filepath.Dir(path), newFileName)

			if err := os.WriteFile(newFilePath, []byte(translatedContent), os.ModePerm); err != nil {
				return fmt.Errorf("failed to write translated file %s: %w", newFilePath, err)
			}
		}

		return nil
	})
}

func cleanBlockDelimiters(content string) string {
	re := regexp.MustCompile("```[a-zA-Z]*\\n|```")
	return re.ReplaceAllString(content, "")
}

func translateText(input, targetLang string, client *openai.Client) (string, error) {
	// Crear un contexto
	ctx := context.Background()

	// Realizar la solicitud de traducción
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "gpt-4o-mini", // Usar el modelo GPT-4o Mini
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a translation model. Translate the following text accurately while maintaining its tone and context.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Translate this text to %s: %s", targetLang, input),
			},
		},
	})
	if err != nil {
		return "", err
	}

	// Validar respuesta
	if len(resp.Choices) == 0 {
		return "", errors.New("no response from OpenAI")
	}

	// Devolver la traducción
	return resp.Choices[0].Message.Content, nil
}
