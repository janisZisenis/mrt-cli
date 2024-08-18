setup() {
    load 'test_helper/common-setup'
    _common_setup
}
@test "can run our script" {
    run mrt

    assert_output 'Hello World'
}
