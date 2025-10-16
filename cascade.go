package main

/*
func tryOpCascade(line string, stack *Stack, ops OpMap) error {
	err := tryOp(line, stack, ops)
	if err != nil {
		return err
	}

	return nil
}
*/

func cascade(line string, stack *Stack, ops *Ops) error {
	err := tryExpr(line, stack)
	if err == nil {
		return nil
	}

	err = ops.Run(line, stack)
	if err == nil {
		return nil
	}

	err = tryCommands(line, stack, ops)
	if err != nil {
		return err
	}

	return nil
}

/*
func tryCascade(line string, stack *Stack, operators OperatorMap) error {
	err := tryExpr(line, stack)
	if err == nil {
		return nil
	}

	err = tryOperator(line, stack, operators)
	if err == nil {
		return nil
	}

	err = tryCommands(line, stack, operators)
	if err != nil {
		return err
	}

	return nil
}
*/
