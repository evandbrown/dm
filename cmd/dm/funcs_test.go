package main

import (
	"testing"
)

/**func TestEnvFunc(t *testing.T) {
	// Create a context to render the template with
	ctx := &Context{
		Data: "env_var = {{env `tpltest`}}",
	}

	// Set an env var
	os.Setenv("tpltest", "val_of_tpltest")
	defer os.Setenv("tpltest", "")

	res, err := Render(ctx)
	if err != nil {
		t.Errorf("parsing: %s", err)
	}
	fmt.Println(res)
}**/

func TestVarFunc(t *testing.T) {
	// Create a context to render the template with
	ctx := &Context{
		Data: "image = {{var `image`}}",
		Vars: map[string]string{
			"image": "image_name",
		},
	}

	_, err := Render(ctx)
	if err != nil {
		t.Errorf("parsing: %s", err)
	}
}
