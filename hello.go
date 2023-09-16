package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 3
const delaySegundos = 5

func main() {
	exibeIntroducao()

	for {
		exibeOpcoes()

		comando := receberComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			imprimeLogs()
		case 0:
			fmt.Println("saindo")
			os.Exit(0)
		default:
			fmt.Println("Erro")
			os.Exit(-1)
		}
	}
}

func exibeIntroducao() {
	fmt.Println("")
	fmt.Println("Ola")
	fmt.Println("Escolha uma das opções abaixo:")
	fmt.Println("")
}

func exibeOpcoes() {
	fmt.Println("1- Iniciar monitoramento")
	fmt.Println("2- Exibir logs")
	fmt.Println("0- Sair do programa")
	fmt.Println("")
}

func receberComando() int {
	var comando int
	fmt.Scan(&comando)
	fmt.Println("")
	return comando
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")
	fmt.Println("")

	sites := leSitesArquivo()

	for i := 0; i < monitoramentos; i++ {
		for position, site := range sites {
			fmt.Println("testando site ", position, ":", site)
			testaSite(site)
		}
		fmt.Println("")
		time.Sleep(delaySegundos * time.Second)
	}
}

func testaSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro", err)
	}

	if resp.StatusCode > 199 && resp.StatusCode < 300 {
		fmt.Println("site: ", site, "carregado com sucesso")
	} else {
		fmt.Println("site: ", site, "falha ao carregar com status code:", resp.StatusCode)
	}

	adicionarLinhaAoLogFile(site, resp.StatusCode)
}

func adicionarLinhaAoLogFile(site string, statusCode int) {
	linha := time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - " + strconv.FormatInt(int64(statusCode), 10)

	arquivo, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Ocorreu um erro :", err)
	}

	arquivo.WriteString(linha)
	arquivo.WriteString("\n")

	arquivo.Close()
}

func leSitesArquivo() []string {
	var sites []string

	arquivo, err := os.Open("SitesMonitorar.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro", err)
	}

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')

		linha = strings.TrimSpace(linha)

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Ocorreu um erro", err)
		}

		sites = append(sites, linha)
	}

	arquivo.Close()

	return sites
}

func imprimeLogs() {
	fmt.Println("exibindo logs...")
	fmt.Println("")

	arquivo, err := os.Open("log.txt")

	if err != nil {

		if strings.Contains(err.Error(), "The system cannot find the file specified") {
			fmt.Println("O arquivo de logs ainda não existe, rode a função 1 do menu para monitorar os sites e gerar o arquivo de log automaticamente")
			return
		}

		fmt.Println(err)
		return
	}

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("error : ", err)
			break
		}

		linha = strings.TrimSpace(linha)

		fmt.Println(linha)
	}

	fmt.Println("")
}
