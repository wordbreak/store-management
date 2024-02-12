

goose:
ifeq ($(c),create)
	goose -dir tools/migrations create $(name) sql
else ifeq ($(env), $(filter $(env),local dev))
	goose -dir tools/migrations mysql store_mgmt_admin:store_mgmt_admin_pass@tcp\(localhost:3306\)/store_mgmt?parseTime=true $(c)
else ifeq ($(env), test)
	goose -dir tools/migrations mysql store_mgmt_admin:store_mgmt_admin_pass@tcp\(localhost:3306\)/store_mgmt_test?parseTime=true $(c)
endif
