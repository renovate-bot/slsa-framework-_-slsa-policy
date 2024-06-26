package deployment

import (
	"os"

	"github.com/slsa-framework/slsa-policy/cli/evaluator/internal/deployment/evaluate"
	"github.com/slsa-framework/slsa-policy/cli/evaluator/internal/deployment/validate"
	"github.com/slsa-framework/slsa-policy/cli/evaluator/internal/utils"
)

func usage(cli string) {
	msg := "" +
		"Usage: %s deployment [options]\n" +
		"\n" +
		"Available options:\n" +
		"validate \t\tValidate the policy files\n" +
		"evaluate \t\tEvaluate the policy\n" +
		"\n"
	utils.Log(msg, cli)
	os.Exit(1)
}

func Run(cli string, args []string) error {
	if len(args) < 1 {
		usage(cli)
	}
	var err error
	switch args[0] {
	default:
		usage(cli)
	case "validate":
		err = validate.Run(cli, args[1:])
	case "evaluate":
		err = evaluate.Run(cli, args[1:])
	}
	return err
}
