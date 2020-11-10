package centralsim

import (
	"log"
	"testing"
)

func TestSimulationEngine1(t *testing.T) {
	log.Printf("************************** Basic simulation test 1 ....")

	//t.Skip("skipping test simulation 1.")
	lfs := Lefs{ // Ejemplo PN datos literales
		IaRed: TransitionList{
			Transition{
				IiIndLocal:        0,
				IiValorLef:        0,
				IiDuracionDisparo: 1,
				TransConstIul:     [][2]int{{0, 1}},
				TransConstPul: [][2]int{
					{1, -1},
					{2, -1},
				},
			},
			Transition{
				IiIndLocal:        1,
				IiValorLef:        1,
				IiDuracionDisparo: 1,
				TransConstIul:     [][2]int{{1, 1}},
				TransConstPul:     [][2]int{{2, -1}},
			},
			Transition{
				IiIndLocal:        2,
				IiValorLef:        2,
				IiDuracionDisparo: 1,
				TransConstIul:     [][2]int{{2, 2}},
				TransConstPul:     [][2]int{{0, -1}},
			},
		},
	}
	ms := MakeMotorSimulation(lfs)
	ms.SimularPeriodo(0, 4) // ciclo 0 hasta ciclo 3
}

func TestSimulationEngine2(t *testing.T) {
	log.Printf("************************** Basic simulation test 2 ....")

	//t.Skip("skipping test simulation 2.")

	// cargamos un fichero de estructura Lef en formato json para centralizado
	lefs, err := Load("testdata/2ramasDe2.rdp.subred0.json")
	if err != nil {
		println("Couln't load th pn file !")
	}

	ms := MakeMotorSimulation(lefs)

	ms.SimularPeriodo(0, 8) // ciclo 0 hasta ciclo 8
}
