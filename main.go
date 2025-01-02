package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func main() {
	foldersEnv := os.Getenv("INPUT_TRANSLATIONS_FOLDER")
	if foldersEnv == "" {
		fmt.Println("Error: INPUT_TRANSLATIONS_FOLDER not set")
		os.Exit(1)
	}

	baseLanguage := os.Getenv("INPUT_BASE_LANGUAGE")
	if baseLanguage == "" {
		fmt.Println("Error: INPUT_BASE_LANGUAGE not set")
		os.Exit(1)
	}

	targetLanguages := os.Getenv("INPUT_TARGET_LANGUAGES")
	if targetLanguages == "" {
		fmt.Println("Error: INPUT_TARGET_LANGUAGES not set")
		os.Exit(1)
	}
	languages := strings.Split(targetLanguages, ",")

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY not set")
		os.Exit(1)
	}
	client := openai.NewClient(apiKey)

	// Dividir rutas proporcionadas en múltiples carpetas
	folders := strings.Split(foldersEnv, ",")

	// Procesar archivos en cada carpeta
	for _, folder := range folders {
		folder = strings.TrimSpace(folder) // Limpiar espacios en blanco
		if err := processFolder(folder, baseLanguage, languages, client); err != nil {
			fmt.Printf("Error processing folder %s: %v\n", folder, err)
		}
	}

	fmt.Println("Translation completed successfully!")
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
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Obtener la extensión y el nombre del archivo
		ext := filepath.Ext(path)
		baseName := strings.TrimSuffix(filepath.Base(path), ext)

		// Traducir y crear archivos para cada idioma
		for _, lang := range languages {
			if lang == baseLanguage {
				continue // Evitar traducir al mismo idioma base
			}

			translatedContent, err := translateText(string(content), lang, client)
			if err != nil {
				return fmt.Errorf("failed to translate file %s to %s: %w", path, lang, err)
			}

			// Crear archivo con el formato <CountryCode>.<Extension>
			newFileName := fmt.Sprintf("%s.%s%s", baseName, lang, ext)
			newFilePath := filepath.Join(filepath.Dir(path), newFileName)

			if err := ioutil.WriteFile(newFilePath, []byte(translatedContent), os.ModePerm); err != nil {
				return fmt.Errorf("failed to write translated file %s: %w", newFilePath, err)
			}
		}

		return nil
	})
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
