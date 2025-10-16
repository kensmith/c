package main

func tryCascade(line string, stack *Stack, operators OperatorMap) error {
	err := tryExpr(line, stack)
	if err == nil {
		return nil
	}

	err = tryOperator(line, stack, operators)
	if err == nil {
		return nil
	}

	err = tryCommands(line, stack)
	if err != nil {
		return err
	}

	return nil
}
