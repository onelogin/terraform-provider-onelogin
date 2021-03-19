resource onelogin_smarthooks basic_test {
  type = "pre-authentication"
  env_vars = [ "API_KEY" ]
  packages = {
    mysql = "2.18.1"
  }
  retries = 0
  timeout = 2
  runtime =  "nodejs12.x"
  disabled = false
  options = {
    risk_enabled = false
    location_enabled = false
  }
  function = <<EOF
		function myFunc() {
			let a = 1;
			let b = 1;
			let c = a + b;
		  console.log("Woo Woo", a, b, c);
		}
	EOF
}
