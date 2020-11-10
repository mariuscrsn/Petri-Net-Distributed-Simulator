package centralsim

import (
	"log"
	"testing"
)

func TestEvent(t *testing.T) {
	log.Printf("************************** Basic event test ....")
	//t.Skip("skipping test evento.")
	log.Printf("Test insercion lista eventos")
	le := EventList{}
	le.inserta(Event{1, 1, 1})
	le.Imprime()
	le.inserta(Event{1, 1, 1})
	le.Imprime()
	le.inserta(Event{0, 3, 3})
	le.Imprime()
}
