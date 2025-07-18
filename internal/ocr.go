package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/genai"
)

const (
	model          = "gemini-2.5-flash"
	gemini_api_key = "[[ YOUR GEMINI API KEY ]]"
	prompt_content = "[[ YOUR PROMPT HERE ]]"
)

func ExtractFromPDF(datadir, outdir string) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  gemini_api_key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(outdir, 0755); err != nil {
		log.Fatalf("Erro ao criar a pasta de saída '%s': %v", outdir, err)
	}
	files, err := os.ReadDir(datadir)
	if err != nil {
		log.Fatalf("Erro ao ler a pasta 'assets': %v", err)
	}

	// 4. Itera sobre cada arquivo na pasta
	for _, file := range files {
		// Ignora subpastas ou arquivos que não sejam PDF
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".pdf") {
			continue
		}

		pdfPath := filepath.Join(datadir, file.Name())
		fmt.Printf("Processando o arquivo: %s\n", pdfPath)

		// Lê o conteúdo do arquivo PDF
		pdfFile, err := os.ReadFile(pdfPath)
		if err != nil {
			log.Printf("Erro ao ler o arquivo '%s': %v", pdfPath, err)
			continue
		}

		// Cria as partes do conteúdo para a requisição
		parts := []*genai.Part{
			&genai.Part{
				InlineData: &genai.Blob{
					MIMEType: "application/pdf",
					Data:     pdfFile,
				},
			},
			genai.NewPartFromText(prompt_content),
		}
		contents := []*genai.Content{
			genai.NewContentFromParts(parts, genai.RoleUser),
		}

		// 5. Envia a requisição para a API do Gemini
		result, err := client.Models.GenerateContent(
			ctx,
			model,
			contents,
			nil,
		)
		if err != nil {
			log.Printf("Erro ao gerar conteúdo para o arquivo '%s': %v", pdfPath, err)
			continue
		}

		// 6. Salva o resultado em um arquivo CSV
		if result.Text() != "" {
			outputFileName := strings.TrimSuffix(file.Name(), ".pdf") + ".csv"
			outputPath := filepath.Join(outdir, outputFileName)

			err = os.WriteFile(outputPath, []byte(result.Text()), 0644)
			if err != nil {
				log.Printf("Erro ao salvar o arquivo '%s': %v", outputPath, err)
			} else {
				fmt.Printf("Dados do arquivo '%s' salvos com sucesso em '%s'.\n", file.Name(), outputPath)
			}
		} else {
			fmt.Printf("A API não retornou texto para o arquivo '%s'.\n", file.Name())
		}

		// Pausa para evitar limites de frequência da API
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\nProcessamento concluído.")
}
