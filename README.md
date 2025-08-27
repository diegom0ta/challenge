# Challenge

### Arquitetura

As duas aplicações vivem no mesmo repositório separadas por camadas como um monólito modular. Na camada ```cmd``` vivem os arquivos de entrada de cada aplicação.

### Prerequisitos

- Docker e Docker Compose instalados
- Make

### Para executar os containers

**Iniciar o banco de dados e rodar as migrações:**
   ```bash
   make up
   ```
## Importante
1. Primeiro deve-se rodar a CLI para fazer a ingestão dos dados. Depois rodar o serviço de API para consulta
2. Os arquivos da B3 são baixados com .txt e é preciso fazer alterações de caracteres para que eles possam ser ingeridos pela CLI:
   ```bash
   mv NOME-DO-ARQUIVO-B3.txt raiz/do/projeto/foo-bar.csv
   ```
   Em seguida fazer as alterações de caracteres via Vim:
   ```bash
   vim nome-do-arquivo.csv
   ```
   No buffer do vim, escrever comando:
   ``` :%! tr "," "." ```
   para substituir virgulas por pontos nos numerais

   ``` :%! tr ";" "," ```
   para substituir ponto-e-vírgulas por vírgulas
   
   Depois salvar o arquivo e sair do Vim   

### Para executar a CLI

Na raiz do projeto:

```bash
go run cmd/cli/main.go "<nome-do-arquivo.csv>"
```
### Para excutar o serviço

Na raiz do projeto:

```bash
go run cmd/api/main.go
```
