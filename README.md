# Cache LRU em Go

## PT-BR

Leia em Ingles: [Readme-En.md](README-En.md)

### Visao geral
Este projeto implementa um cache **LRU (Least Recently Used)** thread-safe em Go.

A estrutura combina:
- um mapa para acesso rapido por chave
- lista duplamente encadeada para manter a ordem de uso (MRU -> LRU)
- `sync.RWMutex` para proteger operacoes concorrentes

Quando a capacidade e excedida, o item menos recentemente usado (LRU) e removido automaticamente.

### Como usar
Pre-requisito: Go instalado.

```go
package main

import (
	"fmt"

	lrucache "github.com/thiagozs/go-lrucache"
)

func main() {
	cache := lrucache.NewLRUCache(lrucache.WithCapacity(2))
	cache.Put(1, 10)
	cache.Put(2, 20)
	cache.Put(3, 30)

	fmt.Println(cache.Get(1)) // -1 (evicted)
	fmt.Println(cache.Get(2)) // 20
	fmt.Println(cache.Get(3)) // 30
}
```

### Rodar testes
```bash
go test ./...
```

### Estrutura do projeto
- [lrucache.go](lrucache.go): implementacao da biblioteca.
- [lrucache_test.go](lrucache_test.go): testes unitarios do pacote.
- [examples/main.go](examples/main.go): exemplo executavel consumindo a biblioteca.


### API principal
- `NewLRUCache(options ...Option) *LRUCache`: cria um cache com capacidade automatica por padrao.
- `WithCapacity(capacity int) Option`: define a capacidade explicitamente quando necessario.
- `Get(key int) int`: retorna o valor e move a chave para MRU; retorna `-1` se nao existir.
- `Put(key int, value int)`: insere/atualiza valor; faz eviction se exceder capacidade.
- `Debug()`: imprime o estado atual em ordem MRU -> LRU.

### Estrutura interna
- `head` e `tail` sao nos sentinela.
- Insercoes e promocoes acontecem logo apos `head` (lado MRU).
- O LRU fica imediatamente antes de `tail`.

Fluxo de escrita (`Put`):
1. Se chave ja existe, atualiza valor e move para frente.
2. Se nao existe, adiciona no frente.
3. Se `len(cache) > capacity`, remove o LRU e deleta do mapa.

Fluxo de leitura (`Get`):
1. Se chave existe, move para frente e retorna valor.
2. Se nao existe, retorna `-1`.

### Complexidade
- `Get`: O(1)
- `Put`: O(1)
- Espaco: O(capacity)

### Concorrencia
As operacoes publicas usam lock:
- `Get` e `Put`: `mu.Lock()` (escrita na ordem LRU)
- `Debug`: `mu.RLock()`

Os helpers internos (`add`, `remove`, `moveToFront`, `removeLRU`) nao aplicam lock proprio; a sincronizacao e responsabilidade do metodo publico chamador.

### Pasta `examples`
- O projeto inclui [examples/main.go](examples/main.go) como exemplo de uso.

### Limitacoes atuais
- Nao ha validacao explicita para `capacity <= 0`.

## Licença

Este projeto é distribuído sob a Licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

## Autor

2026, Thiago Zilli Sarmento :heart: