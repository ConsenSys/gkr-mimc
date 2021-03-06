package polynomial

import (
	"github.com/consensys/gnark/frontend"
)

// UnivariateEqEval computes f(q, h) = 1 - q - h + 2 * q * h
// It returns 1 if q == h \in {0, 1}
func UnivariateEqEval(cs *frontend.ConstraintSystem, q, h frontend.Variable) frontend.Variable {
	res := cs.Constant(1)
	res = cs.Sub(res, q)
	res = cs.Sub(res, h)
	res = cs.Add(res, cs.Mul(cs.Constant(2), cs.Mul(q, h)))
	return res
}

// EqEval returns Eq(q', h')
func EqEval(cs *frontend.ConstraintSystem, qPrime, hPrime []frontend.Variable) frontend.Variable {
	res := cs.Constant(1)
	// multiply all the UnivariateEqEval's into res
	for i := range qPrime {
		res = cs.Mul(res, UnivariateEqEval(cs, qPrime[i], hPrime[i]))
	}
	return res
}
