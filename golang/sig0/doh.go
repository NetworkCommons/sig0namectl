package sig0

import (
	"fmt"
	"io"
	"net/http"

	"github.com/miekg/dns"
)

func SendDOHQuery(server, q string) (*dns.Msg, error) {
	// send over DoH
	url := fmt.Sprintf("https://%s/dns-query?dns=%s", server, q)
	fmt.Println("Q:(DoH):", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Status: %s", resp.Status)
	}

	answerBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var answer = new(dns.Msg)
	err = answer.Unpack(answerBody)
	if err != nil {
		return nil, err
	}

	return answer, nil
}

/*
func SendUDPQuery(server, q string) (*dns.Msg, error) {

   }
*/
