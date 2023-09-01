# Requisitos Aplicacao (Server)
- [x] O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida deverá retornar no formato JSON o resultado para o cliente.
- [x] Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.
- [x] O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e a porta a ser utilizada pelo servidor HTTP será a 8080.
# Requisitos Aplicacao (Client)
- [x] O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.
- [x] O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON). Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.
- [x] O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}  

OBS: Somente os timeouts que não consegui no tempo esperado

**Corrigido o problema da tabela no sqlite

***Mudei para Nanoseconds e realmente funcionou, gostaria de compreender mais essas medidas de tempo, por que confesso que ficou bem confuso pra mim.

****Compreendi, o time.After garante que se atingir os 10 milissegundos ou mais ele retorne com sucesso, tirei o select case e reaproveitei as variáveis ctx e cancel já criadas com o objetivo de setar um novo context time para o banco de dados e usei o execContext para colocar um context especifico, nesse caso, também alterei o context timeout do client para 300 conforme o exercicio.
