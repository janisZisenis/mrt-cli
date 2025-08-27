mrt_execute() {
	bats_load_library 'fixtures/common_fixture.bash'

	"$(mrt run binary-location -- --exe-name)" --team-dir "$(test_env_dir)" "$@"
}
