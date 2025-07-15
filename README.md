# Projeto Desafio Multithreading

O projeto foi desenvolvido utilizando os padrões ensinados no módulo de Multithreading da formação Go Expert - Full Cycle

## Execução

- Rodar comando `go mod tidy`
- Rodar comando `go run main.go`
- Em um novo terminal rodar o comando `curl localhost:8080/obter-cep/89128970`
- Irá retornar um response contendo dados de cep e qual API retornou.

### Retorno Esperado:
```
{
    "cep":"89128970",
    "state":"SC",
    "city":"Luiz Alves",
    "neighborhood":"Vila do Salto",
    "service":"open-cep",
    "apiRetorno":"Brasil API"
}
```
*Obs: O retorno pode variar entre Brasil API e ViaCep*