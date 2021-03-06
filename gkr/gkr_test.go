package gkr

import (
	"gkr-mimc/circuit"
	"gkr-mimc/common"
	"testing"

	"github.com/consensys/gurvy/bn256/fr"
	"github.com/stretchr/testify/assert"
)

// We test for the following circuit
//
//	  a + b + (a×b)	  (a+b) × a × b				  c + d + (c×d)	  (c+d) × c × d
//	________|_______________|_________			________|_______________|_________
// 	|	_________		_________	 |		 	|	_________		_________	 |
//	|	|		|		|		|	 |			|	|		|		|		|	 |
//	|	|	+	|		|	×	|	 |			|	|	+	|		|	×	|	 |
//	|	|_______|		|_______|	 |			|	|_______|		|_______|	 |
//	|________________________________|			|________________________________|
//			|				|							|				|
//		  a + b			  a × b						  c + d			  c × d
//	________|_______________|_________			________|_______________|_________
// 	|	_________		_________	 |		 	|	_________		_________	 |
//	|	|		|		|		|	 |			|	|		|		|		|	 |
//	|	|	+	|		|	×	|	 |			|	|	+	|		|	×	|	 |
//	|	|_______|		|_______|	 |			|	|_______|		|_______|	 |
//	|________________________________|			|________________________________|
//			|				|							|				|
//			a				b							c				d
//
// Note: bN = bG = 1.

// Testcase:
//
//	        5	            6				           19				84
//	________|_______________|_________			________|_______________|_________
// 	|	_________		_________	 |		 	|	_________		_________	 |
//	|	|		|		|		|	 |			|	|		|		|		|	 |
//	|	|	+	|		|	×	|	 |			|	|	+	|		|	×	|	 |
//	|	|_______|		|_______|	 |			|	|_______|		|_______|	 |
//	|________________________________|			|________________________________|
//			|				|							|				|
//		    3			    2						    7			   12
//	________|_______________|_________			________|_______________|_________
// 	|	_________		_________	 |		 	|	_________		_________	 |
//	|	|		|		|		|	 |			|	|		|		|		|	 |
//	|	|	+	|		|	×	|	 |			|	|	+	|		|	×	|	 |
//	|	|_______|		|_______|	 |			|	|_______|		|_______|	 |
//	|________________________________|			|________________________________|
//			|				|							|				|
//			1				2							3				4

func TestMultiBGs(t *testing.T) {
	var one fr.Element
	one.SetOne()

	c := circuit.NewCircuit(
		[][]circuit.Wire{
			// Layer 0
			[]circuit.Wire{
				circuit.Wire{L: 0, R: 1, O: 0, Gate: circuit.AddGate{}},
				circuit.Wire{L: 0, R: 1, O: 1, Gate: circuit.MulGate{}},
			},
			// Layer 1
			[]circuit.Wire{
				circuit.Wire{L: 0, R: 1, O: 0, Gate: circuit.AddGate{}},
			},
		},
	)

	inputs := [][]fr.Element{
		[]fr.Element{common.Uint64ToFr(1), common.Uint64ToFr(2)},
		[]fr.Element{common.Uint64ToFr(3), common.Uint64ToFr(4)},
	}

	a := c.Assign(inputs, 2)
	inputsV := append([][]fr.Element{}, inputs...)
	outputsV := a.Values[2]

	p := NewProver(c, a)
	proof := p.Prove(1)
	v := NewVerifier(1, c)
	validity := v.Verify(proof, inputsV, outputsV)

	assert.Equal(
		t,
		validity,
		true,
		"Proof invalid.",
	)

	actualValues := [][][]fr.Element{
		[][]fr.Element{
			[]fr.Element{common.Uint64ToFr(1), common.Uint64ToFr(2)},
			[]fr.Element{common.Uint64ToFr(3), common.Uint64ToFr(4)},
		},
		[][]fr.Element{
			[]fr.Element{common.Uint64ToFr(3), common.Uint64ToFr(2)},
			[]fr.Element{common.Uint64ToFr(7), common.Uint64ToFr(12)},
		},
		[][]fr.Element{
			[]fr.Element{common.Uint64ToFr(5)},
			[]fr.Element{common.Uint64ToFr(19)},
		},
	}

	assert.Equal(
		t,
		a.Values,
		actualValues,
		"Assignment invalid.",
	)
}

func TestGKR(t *testing.T) {
	var one fr.Element
	one.SetOne()

	c := circuit.NewCircuit(
		[][]circuit.Wire{
			// Layer 0
			[]circuit.Wire{
				circuit.Wire{L: 0, R: 1, O: 0, Gate: circuit.AddGate{}},
				circuit.Wire{L: 0, R: 1, O: 1, Gate: circuit.MulGate{}},
			},
			// Layer 1
			[]circuit.Wire{
				circuit.Wire{L: 0, R: 1, O: 0, Gate: circuit.AddGate{}},
				circuit.Wire{L: 0, R: 1, O: 1, Gate: circuit.MulGate{}},
			},
		},
	)

	inputs := [][]fr.Element{
		[]fr.Element{common.Uint64ToFr(1), common.Uint64ToFr(2)},
		[]fr.Element{common.Uint64ToFr(3), common.Uint64ToFr(4)},
	}

	a := c.Assign(inputs, 2)

	expectedValues := [][][]fr.Element{
		[][]fr.Element{
			[]fr.Element{common.Uint64ToFr(1), common.Uint64ToFr(2)},
			[]fr.Element{common.Uint64ToFr(3), common.Uint64ToFr(4)},
		},
		[][]fr.Element{
			[]fr.Element{common.Uint64ToFr(3), common.Uint64ToFr(2)},
			[]fr.Element{common.Uint64ToFr(7), common.Uint64ToFr(12)},
		},
		[][]fr.Element{
			[]fr.Element{common.Uint64ToFr(5), common.Uint64ToFr(6)},
			[]fr.Element{common.Uint64ToFr(19), common.Uint64ToFr(84)},
		},
	}

	assert.Equal(
		t,
		expectedValues,
		a.Values,
		"Assignment invalid.",
	)

	outputs := a.Values[2]

	p := NewProver(c, a)
	proof := p.Prove(1)
	v := NewVerifier(1, c)
	validity := v.Verify(proof, inputs, outputs)

	assert.Equal(
		t,
		validity,
		true,
		"Proof invalid.",
	)
}
