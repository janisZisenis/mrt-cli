setup() {
    load 'test_helper/bats-support/load'
    load 'test_helper/bats-assert/load'
}

@test "can run our script" {
    run ./build/mrt

    assert_output 'Hello World'
}
