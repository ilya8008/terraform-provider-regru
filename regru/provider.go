package main

import (
        "github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
        return &schema.Provider{
                ResourcesMap: map[string]*schema.Resource{
                        "regru_server": regru.resourceServer(),
                        "regru_ssh": regru.resourceSSH(),
                },
        }
}
