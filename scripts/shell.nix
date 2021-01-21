{ pkgs ? import <nixpkgs> {} }:

with pkgs;

mkShell {
	name="dev-environment";
	buildInputs = [
		go
		go-swag
	];
	shellHook = ''
		shopt -s expand_aliases
		export BROWSER=chromium
		export CHROME_BIN=chromium
		alias build="cd cmd/service; go build; cd ../../"
		alias run="./cmd/service/service start -c configs/appconfig.json"
		alias test="echo 'TODO:'"
		alias updateAPI="sh scripts/GenerateSwaggerDoc.sh"
		alias help="echo 'command: build, run, test, updateAPI'"
		echo "type help to see all available commands"
	'';
}