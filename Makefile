MOD = github.com/Alf-Grindel/clide

.PHONY: gen_model_user
gen_model_user:
	hz model --mod=$(MOD) --idl=idl/user.thrift --model_dir=internal/model

.PHONY: run
run:
	cd cmd && go run main.go